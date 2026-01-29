/**
 * auth.go
 * 認証管理コマンド
 *
 * Ghost Admin APIキーの追加、一覧表示、削除、状態確認を行います。
 */

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mtane0412/gho/internal/config"
	"github.com/mtane0412/gho/internal/ghostapi"
	"github.com/mtane0412/gho/internal/outfmt"
	"github.com/mtane0412/gho/internal/secrets"
)

// AuthCmd は認証管理コマンドのルートです
type AuthCmd struct {
	Add    AuthAddCmd    `cmd:"" help:"Add a new site authentication"`
	List   AuthListCmd   `cmd:"" help:"List authenticated sites"`
	Remove AuthRemoveCmd `cmd:"" help:"Remove site authentication"`
	Status AuthStatusCmd `cmd:"" help:"Check authentication status"`
}

// AuthAddCmd はサイト認証を追加するコマンドです
type AuthAddCmd struct {
	SiteURL string `arg:"" help:"Ghost site URL (e.g., https://myblog.ghost.io)"`
	Alias   string `help:"Alias for this site" short:"a"`
}

// Run はauth addコマンドを実行します
func (c *AuthAddCmd) Run(root *RootFlags) error {
	// 設定ファイルパスを取得
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// 設定を読み込む
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("設定の読み込みに失敗: %w", err)
	}

	// APIキーを入力してもらう
	fmt.Print("Enter Admin API Key (id:secret): ")
	reader := bufio.NewReader(os.Stdin)
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("APIキーの読み込みに失敗: %w", err)
	}
	apiKey = strings.TrimSpace(apiKey)

	// APIキーをパース
	keyID, secret, err := secrets.ParseAdminAPIKey(apiKey)
	if err != nil {
		return err
	}

	// APIキーを検証（サイト情報を取得してみる）
	client, err := ghostapi.NewClient(c.SiteURL, keyID, secret)
	if err != nil {
		return err
	}

	site, err := client.GetSite()
	if err != nil {
		return fmt.Errorf("APIキーの検証に失敗: %w", err)
	}

	// エイリアスを決定
	alias := c.Alias
	if alias == "" {
		// URLからエイリアスを生成（例: https://myblog.ghost.io -> myblog）
		alias = extractAliasFromURL(c.SiteURL)
	}

	// キーリングに保存
	store, err := secrets.NewStore(cfg.KeyringBackend, getKeyringDir())
	if err != nil {
		return fmt.Errorf("キーリングのオープンに失敗: %w", err)
	}

	if err := store.Set(alias, apiKey); err != nil {
		return err
	}

	// 設定にサイトを追加
	cfg.AddSite(alias, c.SiteURL)
	if cfg.DefaultSite == "" {
		cfg.DefaultSite = alias
	}

	// 設定を保存
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("設定の保存に失敗: %w", err)
	}

	fmt.Printf("✓ Added site '%s' (%s)\n", alias, site.Title)
	return nil
}

// AuthListCmd は認証済みサイトを一覧表示するコマンドです
type AuthListCmd struct{}

// Run はauth listコマンドを実行します
func (c *AuthListCmd) Run(root *RootFlags) error {
	// 設定ファイルパスを取得
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// 設定を読み込む
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("設定の読み込みに失敗: %w", err)
	}

	// サイト一覧を取得
	if len(cfg.Sites) == 0 {
		fmt.Println("No sites configured")
		return nil
	}

	// 出力フォーマッターを作成
	formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())

	// テーブル形式で出力
	headers := []string{"Alias", "URL", "Default"}
	var rows [][]string
	for alias, url := range cfg.Sites {
		isDefault := ""
		if alias == cfg.DefaultSite {
			isDefault = "*"
		}
		rows = append(rows, []string{alias, url, isDefault})
	}

	return formatter.PrintTable(headers, rows)
}

// AuthRemoveCmd はサイト認証を削除するコマンドです
type AuthRemoveCmd struct {
	Alias string `arg:"" help:"Site alias to remove"`
}

// Run はauth removeコマンドを実行します
func (c *AuthRemoveCmd) Run(root *RootFlags) error {
	// 設定ファイルパスを取得
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// 設定を読み込む
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("設定の読み込みに失敗: %w", err)
	}

	// サイトが存在するか確認
	if _, ok := cfg.Sites[c.Alias]; !ok {
		return fmt.Errorf("サイトエイリアス '%s' が見つかりません", c.Alias)
	}

	// 確認（--forceフラグがない場合）
	if !root.Force {
		fmt.Printf("Remove site '%s'? (y/N): ", c.Alias)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	// キーリングから削除
	store, err := secrets.NewStore(cfg.KeyringBackend, getKeyringDir())
	if err != nil {
		return fmt.Errorf("キーリングのオープンに失敗: %w", err)
	}

	if err := store.Delete(c.Alias); err != nil {
		return err
	}

	// 設定から削除
	delete(cfg.Sites, c.Alias)
	if cfg.DefaultSite == c.Alias {
		cfg.DefaultSite = ""
	}

	// 設定を保存
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("設定の保存に失敗: %w", err)
	}

	fmt.Printf("✓ Removed site '%s'\n", c.Alias)
	return nil
}

// AuthStatusCmd は認証状態を確認するコマンドです
type AuthStatusCmd struct{}

// Run はauth statusコマンドを実行します
func (c *AuthStatusCmd) Run(root *RootFlags) error {
	// 設定ファイルパスを取得
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// 設定を読み込む
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("設定の読み込みに失敗: %w", err)
	}

	if len(cfg.Sites) == 0 {
		fmt.Println("No sites configured")
		fmt.Println("Run 'gho auth add <site-url>' to add a site")
		return nil
	}

	fmt.Printf("Default site: %s\n", cfg.DefaultSite)
	fmt.Printf("Configured sites: %d\n", len(cfg.Sites))
	fmt.Printf("Keyring backend: %s\n", cfg.KeyringBackend)

	return nil
}

// getConfigPath は設定ファイルのパスを返します
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "gho", "config.json"), nil
}

// getKeyringDir はキーリングディレクトリのパスを返します
func getKeyringDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "gho", "keyring")
}

// extractAliasFromURL はURLからエイリアスを抽出します
func extractAliasFromURL(url string) string {
	// https://myblog.ghost.io -> myblog
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	parts := strings.Split(url, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return "site"
}
