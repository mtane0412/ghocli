/**
 * site.go
 * サイト情報取得コマンド
 *
 * Ghost サイトの基本情報を取得して表示します。
 */

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/config"
	"github.com/mtane0412/gho/internal/errfmt"
	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
	"github.com/mtane0412/gho/internal/secrets"
)

// SiteCmd is the command to retrieve site information
type SiteCmd struct{}

// Run executes the site command
func (c *SiteCmd) Run(ctx context.Context, root *RootFlags) error {
	// Get API client
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// Get site information
	site, err := client.GetSite()
	if err != nil {
		return fmt.Errorf("failed to get site information: %w", err)
	}

	// Create output formatter
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// Output as-is if JSON format
	if root.JSON {
		return formatter.Print(site)
	}

	// Output in key/value format (no headers)
	rows := [][]string{
		{"title", site.Title},
		{"description", site.Description},
		{"url", site.URL},
		{"version", site.Version},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		return err
	}

	return formatter.Flush()
}

// getAPIClient はAPIクライアントを取得します
func getAPIClient(root *RootFlags) (*ghostapi.Client, error) {
	// 設定ファイルパスを取得
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// 設定を読み込む
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// サイトURLを決定
	siteURL := root.Site
	if siteURL == "" {
		siteURL = cfg.DefaultSite
	}
	if siteURL == "" {
		return nil, errors.New(errfmt.FormatSiteError())
	}

	// エイリアスの場合はURLに変換
	if url, ok := cfg.GetSiteURL(siteURL); ok {
		siteURL = url
	} else {
		return nil, fmt.Errorf("site '%s' not found", siteURL)
	}

	// エイリアスを逆引き
	alias := ""
	for a, u := range cfg.Sites {
		if u == siteURL {
			alias = a
			break
		}
	}
	if alias == "" {
		return nil, fmt.Errorf("alias not found for site URL")
	}

	// キーリングからAPIキーを取得
	store, err := secrets.NewStore(cfg.KeyringBackend, getKeyringDir())
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	apiKey, err := store.Get(alias)
	if err != nil {
		return nil, errors.New(errfmt.FormatAuthError(alias))
	}

	// APIキーをパース
	keyID, secret, err := secrets.ParseAdminAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	// APIクライアントを作成
	return ghostapi.NewClient(siteURL, keyID, secret)
}
