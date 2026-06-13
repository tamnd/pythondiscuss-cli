package pythondiscuss_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tamnd/pythondiscuss-cli/pythondiscuss"
)

func newTestClient(ts *httptest.Server) *pythondiscuss.Client {
	cfg := pythondiscuss.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return pythondiscuss.NewClient(cfg)
}

func TestLatest(t *testing.T) {
	payload := map[string]any{
		"topic_list": map[string]any{
			"topics": []map[string]any{
				{
					"id":          1001,
					"title":       "PEP 750: Template strings",
					"created_at":  "2024-01-15T10:30:00.000Z",
					"views":       1500,
					"reply_count": 42,
					"like_count":  38,
					"posts_count": 45,
					"category_id": 3,
					"tags":        []map[string]any{{"id": 1, "name": "pep", "slug": "pep"}, {"id": 2, "name": "syntax", "slug": "syntax"}},
				},
				{
					"id":          1002,
					"title":       "Python 3.14 release schedule",
					"created_at":  "2024-01-14T09:00:00.000Z",
					"views":       820,
					"reply_count": 18,
					"like_count":  20,
					"posts_count": 20,
					"category_id": 5,
					"tags":        []map[string]any{},
				},
			},
		},
		"users": []any{},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/latest.json") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	topics, err := c.Latest(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(topics) != 2 {
		t.Fatalf("got %d topics, want 2", len(topics))
	}
	if topics[0].ID != 1001 {
		t.Errorf("topics[0].ID = %d, want 1001", topics[0].ID)
	}
	if topics[0].Title != "PEP 750: Template strings" {
		t.Errorf("topics[0].Title = %q", topics[0].Title)
	}
	if topics[0].Views != 1500 {
		t.Errorf("topics[0].Views = %d, want 1500", topics[0].Views)
	}
	if topics[0].Replies != 42 {
		t.Errorf("topics[0].Replies = %d, want 42", topics[0].Replies)
	}
	if topics[0].Tags != "pep, syntax" {
		t.Errorf("topics[0].Tags = %q, want %q", topics[0].Tags, "pep, syntax")
	}
	wantURL := "https://discuss.python.org/t/1001"
	if topics[0].URL != wantURL {
		t.Errorf("topics[0].URL = %q, want %q", topics[0].URL, wantURL)
	}
}

func TestTop(t *testing.T) {
	payload := map[string]any{
		"topic_list": map[string]any{
			"topics": []map[string]any{
				{
					"id":          2001,
					"title":       "Top monthly topic",
					"created_at":  "2024-01-01T00:00:00.000Z",
					"views":       5000,
					"reply_count": 100,
					"like_count":  200,
					"posts_count": 102,
					"category_id": 1,
					"tags":        []map[string]any{},
				},
			},
		},
	}
	var capturedQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	topics, err := c.Top(context.Background(), "monthly", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(topics) != 1 {
		t.Fatalf("got %d topics, want 1", len(topics))
	}
	if topics[0].ID != 2001 {
		t.Errorf("topics[0].ID = %d, want 2001", topics[0].ID)
	}
	if !strings.Contains(capturedQuery, "period=monthly") {
		t.Errorf("query %q does not contain period=monthly", capturedQuery)
	}
}

func TestSearch(t *testing.T) {
	payload := map[string]any{
		"posts": []any{},
		"topics": []map[string]any{
			{
				"id":          3001,
				"title":       "typing.Protocol discussion",
				"created_at":  "2024-01-10T12:00:00.000Z",
				"views":       300,
				"reply_count": 15,
				"like_count":  10,
				"posts_count": 17,
				"category_id": 2,
				"tags":        []map[string]any{{"id": 3, "name": "typing", "slug": "typing"}},
			},
		},
	}
	var capturedQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	topics, err := c.Search(context.Background(), "typing protocol", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(topics) != 1 {
		t.Fatalf("got %d topics, want 1", len(topics))
	}
	if topics[0].ID != 3001 {
		t.Errorf("topics[0].ID = %d, want 3001", topics[0].ID)
	}
	if !strings.Contains(capturedQuery, "q=") {
		t.Errorf("query %q does not contain q=", capturedQuery)
	}
}

func TestGetTopic(t *testing.T) {
	payload := map[string]any{
		"id":          4001,
		"title":       "PEP discussion",
		"created_at":  "2024-01-05T08:00:00.000Z",
		"views":       200,
		"reply_count": 5,
		"like_count":  8,
		"posts_count": 7,
		"category_id": 3,
		"tags":        []map[string]any{},
		"post_stream": map[string]any{
			"posts": []map[string]any{
				{
					"id":         50001,
					"topic_id":   4001,
					"username":   "guido",
					"created_at": "2024-01-05T08:00:00.000Z",
					"cooked":     "<p>Template strings provide <em>power</em> &amp; flexibility.</p>",
					"like_count": 15,
				},
				{
					"id":         50002,
					"topic_id":   4001,
					"username":   "brett",
					"created_at": "2024-01-05T09:00:00.000Z",
					"cooked":     "<p>I agree with this approach.</p>",
					"like_count": 5,
				},
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/t/4001.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	posts, err := c.GetTopic(context.Background(), 4001)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 2 {
		t.Fatalf("got %d posts, want 2", len(posts))
	}
	if posts[0].ID != 50001 {
		t.Errorf("posts[0].ID = %d, want 50001", posts[0].ID)
	}
	if posts[0].Username != "guido" {
		t.Errorf("posts[0].Username = %q, want %q", posts[0].Username, "guido")
	}
	// HTML stripped: tags removed, entities decoded
	wantBody := "Template strings provide power & flexibility."
	if posts[0].Body != wantBody {
		t.Errorf("posts[0].Body = %q, want %q", posts[0].Body, wantBody)
	}
	if posts[0].Likes != 15 {
		t.Errorf("posts[0].Likes = %d, want 15", posts[0].Likes)
	}
}

func TestCategories(t *testing.T) {
	payload := map[string]any{
		"category_list": map[string]any{
			"categories": []map[string]any{
				{
					"id":          1,
					"name":        "Python Discussion",
					"description": "<p>Discussion about Python language &amp; ecosystem.</p>",
					"topic_count": 450,
					"post_count":  8200,
					"color":       "0099dd",
				},
				{
					"id":          3,
					"name":        "Packaging",
					"description": "<p>Packaging tools and standards.</p>",
					"topic_count": 320,
					"post_count":  5100,
					"color":       "e45735",
				},
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/categories.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	cats, err := c.Categories(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(cats) != 2 {
		t.Fatalf("got %d categories, want 2", len(cats))
	}
	if cats[0].ID != 1 {
		t.Errorf("cats[0].ID = %d, want 1", cats[0].ID)
	}
	if cats[0].Name != "Python Discussion" {
		t.Errorf("cats[0].Name = %q", cats[0].Name)
	}
	// HTML stripped + entity decoded
	wantDesc := "Discussion about Python language & ecosystem."
	if cats[0].Description != wantDesc {
		t.Errorf("cats[0].Description = %q, want %q", cats[0].Description, wantDesc)
	}
	if cats[0].Topics != 450 {
		t.Errorf("cats[0].Topics = %d, want 450", cats[0].Topics)
	}
	if cats[0].Posts != 8200 {
		t.Errorf("cats[0].Posts = %d, want 8200", cats[0].Posts)
	}
}

func TestGetTopicNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.GetTopic(context.Background(), 99999)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestUserAgent(t *testing.T) {
	var gotUA string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		payload := map[string]any{
			"topic_list": map[string]any{"topics": []any{}},
		}
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.Latest(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if gotUA == "" {
		t.Error("request carried no User-Agent")
	}
}
