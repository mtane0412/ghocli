/**
 * types.go
 * 共通型定義
 *
 * Ghost Admin APIの共通型（Author等）を定義します。
 */

package ghostapi

// Author はGhostの著者を表します
type Author struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Slug     string `json:"slug,omitempty"`
	Email    string `json:"email,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Location string `json:"location,omitempty"`
	Website  string `json:"website,omitempty"`
}
