/**
 * site.go
 * サイト情報取得コマンド
 *
 * Ghost サイトの基本情報を取得して表示します。
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/config"
	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
	"github.com/mtane0412/gho/internal/secrets"
)

// SiteCmd はサイト情報を取得するコマンドです
type SiteCmd struct{}

// Run はsiteコマンドを実行します
func (c *SiteCmd) Run(root *RootFlags) error {
	// APIクライアントを取得
	client, err := getAPIClient(root)
	if err != nil {
		return err
	}

	// サイト情報を取得
	site, err := client.GetSite()
	if err != nil {
		return fmt.Errorf("サイト情報の取得に失敗: %w", err)
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// JSON形式の場合はそのまま出力
	if root.JSON {
		return formatter.Print(site)
	}

	// キー/値形式で出力（ヘッダーなし）
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
		return nil, fmt.Errorf("設定の読み込みに失敗: %w", err)
	}

	// サイトURLを決定
	siteURL := root.Site
	if siteURL == "" {
		siteURL = cfg.DefaultSite
	}
	if siteURL == "" {
		return nil, fmt.Errorf("サイトが指定されていません。-s フラグでサイトを指定するか、デフォルトサイトを設定してください")
	}

	// エイリアスの場合はURLに変換
	if url, ok := cfg.GetSiteURL(siteURL); ok {
		siteURL = url
	} else {
		return nil, fmt.Errorf("サイト '%s' が見つかりません", siteURL)
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
		return nil, fmt.Errorf("サイトURLに対応するエイリアスが見つかりません")
	}

	// キーリングからAPIキーを取得
	store, err := secrets.NewStore(cfg.KeyringBackend, getKeyringDir())
	if err != nil {
		return nil, fmt.Errorf("キーリングのオープンに失敗: %w", err)
	}

	apiKey, err := store.Get(alias)
	if err != nil {
		return nil, fmt.Errorf("APIキーの取得に失敗: %w", err)
	}

	// APIキーをパース
	keyID, secret, err := secrets.ParseAdminAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	// APIクライアントを作成
	return ghostapi.NewClient(siteURL, keyID, secret)
}
