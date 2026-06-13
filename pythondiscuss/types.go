package pythondiscuss

import (
	"fmt"
	"regexp"
	"strings"
)

// ─── exported record types ────────────────────────────────────────────────────

// Topic is the record emitted for latest, top, and search results.
type Topic struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	CreatedAt  string `json:"created_at"`
	Views      int    `json:"views"`
	Replies    int    `json:"replies"`
	Likes      int    `json:"likes"`
	Posts      int    `json:"posts"`
	CategoryID int    `json:"category_id"`
	Tags       string `json:"tags"`
	URL        string `json:"url"`
}

// Post is the record emitted for topic detail.
type Post struct {
	ID        int    `json:"id"`
	TopicID   int    `json:"topic_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	Body      string `json:"body"`
	Likes     int    `json:"likes"`
}

// Category is the record emitted for the categories command.
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Topics      int    `json:"topics"`
	Posts       int    `json:"posts"`
}

// ─── wire types ───────────────────────────────────────────────────────────────

type topicListResp struct {
	TopicList struct {
		Topics        []wireTopic `json:"topics"`
		MoreTopicsURL string      `json:"more_topics_url"`
	} `json:"topic_list"`
}

type wireTopic struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	CreatedAt  string   `json:"created_at"`
	BumpedAt   string   `json:"bumped_at"`
	Views      int      `json:"views"`
	ReplyCount int      `json:"reply_count"`
	LikeCount  int      `json:"like_count"`
	PostsCount int      `json:"posts_count"`
	CategoryID int      `json:"category_id"`
	Tags       []string `json:"tags"`
	Pinned     bool     `json:"pinned"`
	Excerpt    string   `json:"excerpt"`
}

type topicDetailResp struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	CreatedAt  string   `json:"created_at"`
	Views      int      `json:"views"`
	ReplyCount int      `json:"reply_count"`
	LikeCount  int      `json:"like_count"`
	PostsCount int      `json:"posts_count"`
	CategoryID int      `json:"category_id"`
	Tags       []string `json:"tags"`
	PostStream struct {
		Posts []wirePost `json:"posts"`
	} `json:"post_stream"`
}

type wirePost struct {
	ID         int    `json:"id"`
	TopicID    int    `json:"topic_id"`
	Username   string `json:"username"`
	CreatedAt  string `json:"created_at"`
	Cooked     string `json:"cooked"`
	LikeCount  int    `json:"like_count"`
	ReplyCount int    `json:"reply_count"`
}

type categoriesResp struct {
	CategoryList struct {
		Categories []wireCategory `json:"categories"`
	} `json:"category_list"`
}

type wireCategory struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TopicCount  int    `json:"topic_count"`
	PostCount   int    `json:"post_count"`
	Color       string `json:"color"`
}

type searchResp struct {
	Posts  []wirePost  `json:"posts"`
	Topics []wireTopic `json:"topics"`
}

// ─── helpers ──────────────────────────────────────────────────────────────────

var reHTMLTag = regexp.MustCompile(`<[^>]+>`)

func stripHTML(s string) string {
	s = reHTMLTag.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	return strings.TrimSpace(s)
}

func topicURL(id int) string {
	return fmt.Sprintf("https://discuss.python.org/t/%d", id)
}

func wireTopicToTopic(w wireTopic) Topic {
	tags := strings.Join(w.Tags, ", ")
	return Topic{
		ID:         w.ID,
		Title:      w.Title,
		CreatedAt:  w.CreatedAt,
		Views:      w.Views,
		Replies:    w.ReplyCount,
		Likes:      w.LikeCount,
		Posts:      w.PostsCount,
		CategoryID: w.CategoryID,
		Tags:       tags,
		URL:        topicURL(w.ID),
	}
}

func wirePostToPost(w wirePost) Post {
	return Post{
		ID:        w.ID,
		TopicID:   w.TopicID,
		Username:  w.Username,
		CreatedAt: w.CreatedAt,
		Body:      stripHTML(w.Cooked),
		Likes:     w.LikeCount,
	}
}

func wireCategoryToCategory(w wireCategory) Category {
	return Category{
		ID:          w.ID,
		Name:        w.Name,
		Description: stripHTML(w.Description),
		Topics:      w.TopicCount,
		Posts:       w.PostCount,
	}
}
