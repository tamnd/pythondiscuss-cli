// Package pythondiscuss is the library behind the pd command: the HTTP client,
// request shaping, and the typed data models for discuss.python.org.
//
// discuss.python.org is a public Discourse forum. All public content is
// available as JSON by appending .json to page URLs or using the documented
// API paths. No API key or authentication is required for reading.
package pythondiscuss

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// DefaultUserAgent identifies the client to discuss.python.org.
const DefaultUserAgent = "pd/dev (+https://github.com/tamnd/pythondiscuss-cli)"

// ErrNotFound is returned when a topic ID does not exist (404 from Discourse).
var ErrNotFound = errors.New("not found")

// Config holds constructor parameters.
type Config struct {
	BaseURL   string
	UserAgent string
	Rate      time.Duration
	Retries   int
	Timeout   time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://discuss.python.org",
		UserAgent: DefaultUserAgent,
		Rate:      500 * time.Millisecond,
		Retries:   3,
		Timeout:   15 * time.Second,
	}
}

// Client talks to discuss.python.org over HTTP.
type Client struct {
	httpClient *http.Client
	cfg        Config
	mu         sync.Mutex
	last       time.Time
}

// NewClient returns a Client with the given config.
func NewClient(cfg Config) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: cfg.Timeout},
		cfg:        cfg,
	}
}

// Latest returns the most recently bumped topics from page N (0-indexed).
func (c *Client) Latest(ctx context.Context, page int) ([]Topic, error) {
	u := fmt.Sprintf("%s/latest.json?page=%d", c.cfg.BaseURL, page)
	var resp topicListResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	topics := make([]Topic, len(resp.TopicList.Topics))
	for i, wt := range resp.TopicList.Topics {
		topics[i] = wireTopicToTopic(wt)
	}
	return topics, nil
}

// Top returns top-ranked topics for the given period (all/yearly/monthly/weekly/daily).
func (c *Client) Top(ctx context.Context, period string, page int) ([]Topic, error) {
	u := fmt.Sprintf("%s/top.json?period=%s&page=%d", c.cfg.BaseURL, url.QueryEscape(period), page)
	var resp topicListResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	topics := make([]Topic, len(resp.TopicList.Topics))
	for i, wt := range resp.TopicList.Topics {
		topics[i] = wireTopicToTopic(wt)
	}
	return topics, nil
}

// Search performs a full-text search and returns matching topics.
// page is 1-indexed (Discourse search pagination).
func (c *Client) Search(ctx context.Context, query string, page int) ([]Topic, error) {
	u := fmt.Sprintf("%s/search.json?q=%s&page=%d", c.cfg.BaseURL, url.QueryEscape(query), page)
	var resp searchResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	topics := make([]Topic, len(resp.Topics))
	for i, wt := range resp.Topics {
		topics[i] = wireTopicToTopic(wt)
	}
	return topics, nil
}

// GetTopic fetches a topic and returns its posts (first page, up to ~20).
func (c *Client) GetTopic(ctx context.Context, id int) ([]Post, error) {
	u := fmt.Sprintf("%s/t/%d.json", c.cfg.BaseURL, id)
	var resp topicDetailResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	posts := make([]Post, len(resp.PostStream.Posts))
	for i, wp := range resp.PostStream.Posts {
		posts[i] = wirePostToPost(wp)
	}
	return posts, nil
}

// Categories returns all forum categories.
func (c *Client) Categories(ctx context.Context) ([]Category, error) {
	u := fmt.Sprintf("%s/categories.json", c.cfg.BaseURL)
	var resp categoriesResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	cats := make([]Category, len(resp.CategoryList.Categories))
	for i, wc := range resp.CategoryList.Categories {
		cats[i] = wireCategoryToCategory(wc)
	}
	return cats, nil
}

// ─── HTTP mechanics ───────────────────────────────────────────────────────────

func (c *Client) getJSON(ctx context.Context, rawURL string, v any) error {
	body, err := c.get(ctx, rawURL)
	if err != nil {
		return err
	}
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "null" {
		return ErrNotFound
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("decode %s: %w", rawURL, err)
	}
	return nil
}

func (c *Client) get(ctx context.Context, rawURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		body, retry, err := c.do(ctx, rawURL)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("get %s: %w", rawURL, lastErr)
}

func (c *Client) do(ctx context.Context, rawURL string) ([]byte, bool, error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, false, ErrNotFound
	}
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, true, err
	}
	return b, false, nil
}

func (c *Client) pace() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 500 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}
