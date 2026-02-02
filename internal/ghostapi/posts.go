/**
 * posts.go
 * Posts API
 *
 * Provides Posts functionality for the Ghost Admin API.
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Post represents a Ghost post
type Post struct {
	// Basic information
	ID     string `json:"id,omitempty"`
	UUID   string `json:"uuid,omitempty"`
	Title  string `json:"title"`
	Slug   string `json:"slug,omitempty"`
	Status string `json:"status"` // draft, published, scheduled
	URL    string `json:"url,omitempty"`

	// Content
	HTML          string `json:"html,omitempty"`
	Lexical       string `json:"lexical,omitempty"`
	Excerpt       string `json:"excerpt,omitempty"`
	CustomExcerpt string `json:"custom_excerpt,omitempty"`

	// Images
	FeatureImage        string `json:"feature_image,omitempty"`
	FeatureImageAlt     string `json:"feature_image_alt,omitempty"`
	FeatureImageCaption string `json:"feature_image_caption,omitempty"`
	OGImage             string `json:"og_image,omitempty"`
	TwitterImage        string `json:"twitter_image,omitempty"`

	// SEO
	MetaTitle          string `json:"meta_title,omitempty"`
	MetaDescription    string `json:"meta_description,omitempty"`
	OGTitle            string `json:"og_title,omitempty"`
	OGDescription      string `json:"og_description,omitempty"`
	TwitterTitle       string `json:"twitter_title,omitempty"`
	TwitterDescription string `json:"twitter_description,omitempty"`
	CanonicalURL       string `json:"canonical_url,omitempty"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`

	// Control
	Visibility string `json:"visibility,omitempty"` // public, members, paid
	Featured   bool   `json:"featured,omitempty"`
	EmailOnly  bool   `json:"email_only,omitempty"`

	// Custom
	CodeinjectionHead string `json:"codeinjection_head,omitempty"`
	CodeinjectionFoot string `json:"codeinjection_foot,omitempty"`
	CustomTemplate    string `json:"custom_template,omitempty"`

	// Related
	Tags          []Tag    `json:"tags,omitempty"`
	Authors       []Author `json:"authors,omitempty"`
	PrimaryAuthor *Author  `json:"primary_author,omitempty"`
	PrimaryTag    *Tag     `json:"primary_tag,omitempty"`

	// Other
	CommentID   string `json:"comment_id,omitempty"`
	ReadingTime int    `json:"reading_time,omitempty"`

	// Email/Newsletter
	EmailSegment           string `json:"email_segment,omitempty"`
	NewsletterID           string `json:"newsletter_id,omitempty"`
	SendEmailWhenPublished bool   `json:"send_email_when_published,omitempty"`
}

// ListOptions represents options for fetching post list
type ListOptions struct {
	Status  string // draft, published, scheduled, all
	Limit   int    // Number of items to fetch (default: 15)
	Page    int    // Page number (default: 1)
	Include string // Additional information to include (tags, authors, etc.)
}

// CreateOptions contains options for creating/updating posts
type CreateOptions struct {
	Source string // "html" for server-side HTML-to-Lexical conversion
}

// PostListResponse represents a post list response
type PostListResponse struct {
	Posts []Post `json:"posts"`
	Meta  struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// ListPosts retrieves a list of posts
func (c *Client) ListPosts(opts ListOptions) (*PostListResponse, error) {
	path := "/ghost/api/admin/posts/"

	// Build query parameters
	params := []string{}
	if opts.Status != "" && opts.Status != "all" {
		params = append(params, fmt.Sprintf("filter=status:%s", opts.Status))
	}
	if opts.Limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", opts.Limit))
	}
	if opts.Page > 0 {
		params = append(params, fmt.Sprintf("page=%d", opts.Page))
	}
	if opts.Include != "" {
		params = append(params, fmt.Sprintf("include=%s", opts.Include))
	}

	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response PostListResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetPost retrieves a post (by ID or slug)
func (c *Client) GetPost(idOrSlug string) (*Post, error) {
	// Determine if it's a slug (IDs are typically 24-character hex strings)
	var path string
	if len(idOrSlug) == 24 {
		// Treat as ID (fetch both HTML and Lexical with formats=html,lexical)
		path = fmt.Sprintf("/ghost/api/admin/posts/%s/?formats=html,lexical", idOrSlug)
	} else {
		// Treat as slug (fetch both HTML and Lexical with formats=html,lexical)
		path = fmt.Sprintf("/ghost/api/admin/posts/slug/%s/?formats=html,lexical", idOrSlug)
	}

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response struct {
		Posts []Post `json:"posts"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Posts) == 0 {
		return nil, fmt.Errorf("post not found: %s", idOrSlug)
	}

	return &response.Posts[0], nil
}

// CreatePost creates a new post
func (c *Client) CreatePost(post *Post) (*Post, error) {
	return c.CreatePostWithOptions(post, CreateOptions{})
}

// CreatePostWithOptions creates a new post with options
func (c *Client) CreatePostWithOptions(post *Post, opts CreateOptions) (*Post, error) {
	path := "/ghost/api/admin/posts/"

	// Create request body
	reqBody := map[string]interface{}{
		"posts": []interface{}{post},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}

	// Build request options
	var reqOpts *RequestOptions
	if opts.Source != "" {
		reqOpts = &RequestOptions{
			QueryParams: map[string]string{
				"source": opts.Source,
			},
		}
	}

	// Execute request
	respBody, err := c.doRequestWithOptions("POST", path, bytes.NewReader(jsonData), reqOpts)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response struct {
		Posts []Post `json:"posts"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Posts) == 0 {
		return nil, fmt.Errorf("failed to create post")
	}

	return &response.Posts[0], nil
}

// UpdatePost updates a post
func (c *Client) UpdatePost(id string, post *Post) (*Post, error) {
	return c.UpdatePostWithOptions(id, post, CreateOptions{})
}

// UpdatePostWithOptions updates a post with options
func (c *Client) UpdatePostWithOptions(id string, post *Post, opts CreateOptions) (*Post, error) {
	path := fmt.Sprintf("/ghost/api/admin/posts/%s/", id)

	// Create request body
	reqBody := map[string]interface{}{
		"posts": []interface{}{post},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}

	// Build request options
	var reqOpts *RequestOptions
	if opts.Source != "" {
		reqOpts = &RequestOptions{
			QueryParams: map[string]string{
				"source": opts.Source,
			},
		}
	}

	// Execute request
	respBody, err := c.doRequestWithOptions("PUT", path, bytes.NewReader(jsonData), reqOpts)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response struct {
		Posts []Post `json:"posts"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Posts) == 0 {
		return nil, fmt.Errorf("failed to update post")
	}

	return &response.Posts[0], nil
}

// DeletePost deletes a post
func (c *Client) DeletePost(id string) error {
	path := fmt.Sprintf("/ghost/api/admin/posts/%s/", id)

	// Execute request
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
