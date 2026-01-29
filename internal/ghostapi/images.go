/**
 * images.go
 * Images API
 *
 * Ghost Admin APIのImages機能を提供します。
 */

package ghostapi

import (
	"encoding/json"
	"fmt"
	"io"
)

// Image はGhostの画像を表します
type Image struct {
	URL string `json:"url"`
	Ref string `json:"ref,omitempty"`
}

// ImageUploadOptions は画像アップロードのオプションです
type ImageUploadOptions struct {
	Purpose string // image, profile_image, icon
	Ref     string // 画像の参照ID
}

// ImageResponse は画像のレスポンスです
type ImageResponse struct {
	Images []Image `json:"images"`
}

// UploadImage は画像をアップロードします
func (c *Client) UploadImage(file io.Reader, filename string, opts ImageUploadOptions) (*Image, error) {
	path := "/ghost/api/admin/images/upload/"

	// マルチパートフィールドを構築
	fields := make(map[string]string)
	if opts.Purpose != "" {
		fields["purpose"] = opts.Purpose
	}
	if opts.Ref != "" {
		fields["ref"] = opts.Ref
	}

	// リクエストを実行
	respBody, err := c.doMultipartRequest(path, file, filename, fields)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	var resp ImageResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("レスポンスのパースに失敗しました: %w", err)
	}

	if len(resp.Images) == 0 {
		return nil, fmt.Errorf("画像のアップロードに失敗しました")
	}

	return &resp.Images[0], nil
}
