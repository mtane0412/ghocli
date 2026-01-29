/**
 * config.go
 * ghoの設定ファイル管理
 *
 * 設定ファイルは ~/.config/gho/config.json に保存され、
 * マルチサイト対応のためのエイリアス機能を提供します。
 */

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Config はghoの設定を表します
type Config struct {
	// KeyringBackend はキーリングのバックエンド種別（auto/file/keychain等）
	KeyringBackend string `json:"keyring_backend"`

	// DefaultSite はデフォルトのサイトエイリアス
	DefaultSite string `json:"default_site,omitempty"`

	// Sites はエイリアスからサイトURLへのマッピング
	Sites map[string]string `json:"sites"`
}

// Load は指定されたパスから設定ファイルを読み込みます。
// ファイルが存在しない場合は、デフォルト値を持つ新しい設定を返します。
func Load(path string) (*Config, error) {
	// ファイルが存在しない場合は、デフォルト設定を返す
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{
			KeyringBackend: "auto",
			DefaultSite:    "",
			Sites:          make(map[string]string),
		}, nil
	}

	// ファイルを読み込む
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// JSONをパース
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Sitesがnilの場合は初期化
	if cfg.Sites == nil {
		cfg.Sites = make(map[string]string)
	}

	return &cfg, nil
}

// Save は設定を指定されたパスに保存します。
func (c *Config) Save(path string) error {
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// JSONに変換（インデント付き）
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	// ファイルに書き込む（0600 = 所有者のみ読み書き可能）
	return os.WriteFile(path, data, 0600)
}

// AddSite はサイトエイリアスとURLを設定に追加します。
func (c *Config) AddSite(alias, url string) {
	if c.Sites == nil {
		c.Sites = make(map[string]string)
	}
	c.Sites[alias] = url
}

// GetSiteURL はエイリアスまたはURL文字列からサイトURLを取得します。
// エイリアスとして登録されている場合は対応するURLを、
// そうでない場合はURL文字列として扱い、そのまま返します。
func (c *Config) GetSiteURL(aliasOrURL string) (string, bool) {
	// エイリアスとして登録されているか確認
	if url, ok := c.Sites[aliasOrURL]; ok {
		return url, true
	}

	// URL文字列として扱う（https://で始まる場合）
	if strings.HasPrefix(aliasOrURL, "https://") || strings.HasPrefix(aliasOrURL, "http://") {
		return aliasOrURL, true
	}

	return "", false
}
