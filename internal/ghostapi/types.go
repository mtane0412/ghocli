/**
 * types.go
 * Common type definitions
 *
 * Defines common types for Ghost Admin API (Author, etc.).
 */

package ghostapi

// Author represents a Ghost author
type Author struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Slug     string `json:"slug,omitempty"`
	Email    string `json:"email,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Location string `json:"location,omitempty"`
	Website  string `json:"website,omitempty"`
}
