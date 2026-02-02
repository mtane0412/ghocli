/**
 * pages.go
 * Pages API
 *
 * Provides Pages functionality for the Ghost Admin API.
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Page represents a Ghost page
type Page struct {
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

// PageListResponse represents a page list response
type PageListResponse struct{
	Pages []Page `json:"pages"`
	Meta  struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

// ListPages retrieves a list of pages
func (c *Client) ListPages(opts ListOptions) (*PageListResponse, error) {
	path := "/ghost/api/admin/pages/"

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
	var response PageListResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetPage retrieves a page (by ID or slug)
func (c *Client) GetPage(idOrSlug string) (*Page, error) {
	// Determine if it's a slug (IDs are typically 24-character hex strings)
	var path string
	if len(idOrSlug) == 24 {
		// Treat as ID (fetch both HTML and Lexical with formats=html,lexical)
		path = fmt.Sprintf("/ghost/api/admin/pages/%s/?formats=html,lexical", idOrSlug)
	} else {
		// Treat as slug (fetch both HTML and Lexical with formats=html,lexical)
		path = fmt.Sprintf("/ghost/api/admin/pages/slug/%s/?formats=html,lexical", idOrSlug)
	}

	// Execute request
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response struct {
		Pages []Page `json:"pages"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Pages) == 0 {
		return nil, fmt.Errorf("page not found: %s", idOrSlug)
	}

	return &response.Pages[0], nil
}

// CreatePage creates a new page
func (c *Client) CreatePage(page *Page) (*Page, error) {
	path := "/ghost/api/admin/pages/"

	// Create request body
	reqBody := map[string]interface{}{
		"pages": []interface{}{page},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}

	// Execute request
	respBody, err := c.doRequest("POST", path, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	// Parse response
	var response struct {
		Pages []Page `json:"pages"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Pages) == 0 {
		return nil, fmt.Errorf("failed to create page")
	}

	return &response.Pages[0], nil
}

// UpdatePage updates a page
func (c *Client) UpdatePage(id string, page *Page) (*Page, error) {
	path := fmt.Sprintf("/ghost/api/admin/pages/%s/", id)

	// Create request body
	reqBody := map[string]interface{}{
		"pages": []interface{}{page},
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}

	// Execute request
	respBody, err := c.doRequest("PUT", path, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	// Parse response
	var response struct {
		Pages []Page `json:"pages"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Pages) == 0 {
		return nil, fmt.Errorf("failed to update page")
	}

	return &response.Pages[0], nil
}

// DeletePage deletes a page
func (c *Client) DeletePage(id string) error {
	path := fmt.Sprintf("/ghost/api/admin/pages/%s/", id)

	// Execute request
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
