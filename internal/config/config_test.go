/**
 * config_test.go
 * 設定システムのテストコード
 */

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadConfig_新規作成時のデフォルト値
func TestLoadConfig_新規作成時のデフォルト値(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// 設定ファイルが存在しない状態でLoadを呼び出す
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("設定ファイルのロードに失敗: %v", err)
	}

	// デフォルト値の検証
	if cfg.KeyringBackend != "auto" {
		t.Errorf("KeyringBackend = %q; want %q", cfg.KeyringBackend, "auto")
	}
	if cfg.DefaultSite != "" {
		t.Errorf("DefaultSite = %q; want empty string", cfg.DefaultSite)
	}
	if cfg.Sites == nil {
		t.Error("Sites map is nil; want empty map")
	}
}

// TestLoadConfig_既存ファイルの読み込み
func TestLoadConfig_既存ファイルの読み込み(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// テスト用の設定ファイルを作成
	configContent := `{
  "keyring_backend": "file",
  "default_site": "myblog",
  "sites": {
    "myblog": "https://myblog.ghost.io",
    "company": "https://blog.company.com"
  }
}`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("テスト用設定ファイルの作成に失敗: %v", err)
	}

	// 設定ファイルを読み込む
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("設定ファイルのロードに失敗: %v", err)
	}

	// 読み込んだ値を検証
	if cfg.KeyringBackend != "file" {
		t.Errorf("KeyringBackend = %q; want %q", cfg.KeyringBackend, "file")
	}
	if cfg.DefaultSite != "myblog" {
		t.Errorf("DefaultSite = %q; want %q", cfg.DefaultSite, "myblog")
	}
	if len(cfg.Sites) != 2 {
		t.Errorf("Sites has %d entries; want 2", len(cfg.Sites))
	}
	if cfg.Sites["myblog"] != "https://myblog.ghost.io" {
		t.Errorf("Sites[myblog] = %q; want %q", cfg.Sites["myblog"], "https://myblog.ghost.io")
	}
	if cfg.Sites["company"] != "https://blog.company.com" {
		t.Errorf("Sites[company] = %q; want %q", cfg.Sites["company"], "https://blog.company.com")
	}
}

// TestSave_設定ファイルの保存
func TestSave_設定ファイルの保存(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// 設定を作成
	cfg := &Config{
		KeyringBackend: "auto",
		DefaultSite:    "testsite",
		Sites: map[string]string{
			"testsite": "https://test.ghost.io",
		},
	}

	// 設定を保存
	if err := cfg.Save(configPath); err != nil {
		t.Fatalf("設定ファイルの保存に失敗: %v", err)
	}

	// ファイルが作成されたことを確認
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("設定ファイルが作成されていない")
	}

	// 保存された設定を再度読み込んで検証
	reloaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("保存した設定ファイルの再読み込みに失敗: %v", err)
	}

	if reloaded.KeyringBackend != cfg.KeyringBackend {
		t.Errorf("KeyringBackend = %q; want %q", reloaded.KeyringBackend, cfg.KeyringBackend)
	}
	if reloaded.DefaultSite != cfg.DefaultSite {
		t.Errorf("DefaultSite = %q; want %q", reloaded.DefaultSite, cfg.DefaultSite)
	}
	if reloaded.Sites["testsite"] != cfg.Sites["testsite"] {
		t.Errorf("Sites[testsite] = %q; want %q", reloaded.Sites["testsite"], cfg.Sites["testsite"])
	}
}

// TestAddSite_サイトの追加
func TestAddSite_サイトの追加(t *testing.T) {
	cfg := &Config{
		KeyringBackend: "auto",
		Sites:          make(map[string]string),
	}

	// サイトを追加
	cfg.AddSite("myblog", "https://myblog.ghost.io")

	// 追加されたサイトを検証
	if cfg.Sites["myblog"] != "https://myblog.ghost.io" {
		t.Errorf("Sites[myblog] = %q; want %q", cfg.Sites["myblog"], "https://myblog.ghost.io")
	}
}

// TestGetSiteURL_エイリアスからURLを取得
func TestGetSiteURL_エイリアスからURLを取得(t *testing.T) {
	cfg := &Config{
		Sites: map[string]string{
			"myblog": "https://myblog.ghost.io",
		},
	}

	// エイリアスからURLを取得
	url, ok := cfg.GetSiteURL("myblog")
	if !ok {
		t.Fatal("GetSiteURL returned false; want true")
	}
	if url != "https://myblog.ghost.io" {
		t.Errorf("url = %q; want %q", url, "https://myblog.ghost.io")
	}

	// 存在しないエイリアス
	_, ok = cfg.GetSiteURL("nonexistent")
	if ok {
		t.Error("GetSiteURL returned true for nonexistent alias; want false")
	}
}

// TestGetSiteURL_URL直接指定
func TestGetSiteURL_URL直接指定(t *testing.T) {
	cfg := &Config{
		Sites: make(map[string]string),
	}

	// URL直接指定（エイリアスとして登録されていない場合はそのまま返す）
	url, ok := cfg.GetSiteURL("https://direct.ghost.io")
	if !ok {
		t.Fatal("GetSiteURL returned false for direct URL; want true")
	}
	if url != "https://direct.ghost.io" {
		t.Errorf("url = %q; want %q", url, "https://direct.ghost.io")
	}
}
