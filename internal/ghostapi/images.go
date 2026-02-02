/**
 * images.go
 * Images API
 *
 * Provides Images functionality for the Ghost Admin API.
 */

package ghostapi

import (
	"encoding/json"
	"fmt"
	"io"
)

// Image represents a Ghost image
type Image struct {
	URL string `json:"url"`
	Ref string `json:"ref,omitempty"`
}

// ImageUploadOptions represents options for image upload
type ImageUploadOptions struct {
	Purpose string // image, profile_image, icon
	Ref     string // Reference ID for the image
}

// ImageResponse represents an image response
type ImageResponse struct {
	Images []Image `json:"images"`
}

// UploadImage uploads an image
func (c *Client) UploadImage(file io.Reader, filename string, opts ImageUploadOptions) (*Image, error) {
	path := "/ghost/api/admin/images/upload/"

	// Build multipart fields
	fields := make(map[string]string)
	if opts.Purpose != "" {
		fields["purpose"] = opts.Purpose
	}
	if opts.Ref != "" {
		fields["ref"] = opts.Ref
	}

	// Execute request
	respBody, err := c.doMultipartRequest(path, file, filename, fields)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp ImageResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(resp.Images) == 0 {
		return nil, fmt.Errorf("failed to upload image")
	}

	return &resp.Images[0], nil
}
