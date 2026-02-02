/**
 * client.go
 * HTTP client for Ghost Admin API
 *
 * Manages HTTP requests to the Ghost Admin API.
 * Each request includes an Authorization header containing a JWT token.
 */

package ghostapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// Client is the Ghost Admin API client
type Client struct {
	baseURL    string
	keyID      string
	secret     string
	httpClient *http.Client
}

// Site represents Ghost site information
type Site struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Version     string `json:"version"`
}

// ErrorResponse represents a Ghost API error response
type ErrorResponse struct {
	Errors []struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"errors"`
}

// NewClient creates a new Ghost Admin API client.
func NewClient(baseURL, keyID, secret string) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("site URL is empty")
	}
	if keyID == "" {
		return nil, errors.New("key ID is empty")
	}
	if secret == "" {
		return nil, errors.New("secret is empty")
	}

	// Remove trailing slash from URL
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &Client{
		baseURL: baseURL,
		keyID:   keyID,
		secret:  secret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// doRequest executes an HTTP request and returns the response body.
func (c *Client) doRequest(method, path string, body io.Reader) ([]byte, error) {
	// Generate JWT token
	token, err := GenerateJWT(c.keyID, c.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Build request URL
	url := c.baseURL + path

	// Create HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Ghost "+token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Parse error response
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && len(errResp.Errors) > 0 {
			return nil, fmt.Errorf("API error: %s", errResp.Errors[0].Message)
		}
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	return respBody, nil
}

// doMultipartRequest executes an HTTP request with multipart/form-data and returns the response body.
func (c *Client) doMultipartRequest(path string, file io.Reader, filename string, fields map[string]string) ([]byte, error) {
	// Generate JWT token
	token, err := GenerateJWT(c.keyID, c.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	// Build multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file field
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create file field: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Add additional fields
	for key, val := range fields {
		if err := writer.WriteField(key, val); err != nil {
			return nil, fmt.Errorf("failed to add field %s: %w", key, err)
		}
	}

	// Close multipart writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Build request URL
	url := c.baseURL + path

	// Create HTTP request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Ghost "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Parse error response
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && len(errResp.Errors) > 0 {
			return nil, fmt.Errorf("API error: %s", errResp.Errors[0].Message)
		}
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	return respBody, nil
}

// GetSite retrieves site information.
func (c *Client) GetSite() (*Site, error) {
	respBody, err := c.doRequest("GET", "/ghost/api/admin/site/", nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response struct {
		Site Site `json:"site"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response.Site, nil
}
