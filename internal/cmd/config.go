/**
 * config.go
 * 設定管理コマンド
 *
 * ghoの設定を管理するコマンド群を提供します。
 */

package cmd

import (
	"fmt"

	"github.com/mtane0412/gho/internal/config"
)

// ConfigCmd は設定管理のルートコマンドです
type ConfigCmd struct {
	Get    ConfigGetCmd    `cmd:"" help:"Get configuration value"`
	Set    ConfigSetCmd    `cmd:"" help:"Set configuration value"`
	Unset  ConfigUnsetCmd  `cmd:"" help:"Unset configuration value"`
	List   ConfigListCmd   `cmd:"" help:"List all configuration"`
	Path   ConfigPathCmd   `cmd:"" help:"Show configuration file path"`
	Keys   ConfigKeysCmd   `cmd:"" help:"List available configuration keys"`
}

// ConfigGetCmd は設定値を取得するコマンドです
type ConfigGetCmd struct {
	Key string `arg:"" help:"Configuration key to get"`
}

// Run はconfig getコマンドを実行します
func (c *ConfigGetCmd) Run(root *RootFlags) error {
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

	// キーに応じて値を取得
	var value string
	switch c.Key {
	case "default_site":
		value = cfg.DefaultSite
	case "keyring_backend":
		value = cfg.KeyringBackend
	default:
		return fmt.Errorf("不明な設定キー: %s", c.Key)
	}

	// 値を出力
	fmt.Println(value)
	return nil
}

// ConfigSetCmd は設定値を設定するコマンドです
type ConfigSetCmd struct {
	Key   string `arg:"" help:"Configuration key to set"`
	Value string `arg:"" help:"Configuration value to set"`
}

// Run はconfig setコマンドを実行します
func (c *ConfigSetCmd) Run(root *RootFlags) error {
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

	// キーに応じて値を設定
	switch c.Key {
	case "default_site":
		cfg.DefaultSite = c.Value
	case "keyring_backend":
		cfg.KeyringBackend = c.Value
	default:
		return fmt.Errorf("不明な設定キー: %s", c.Key)
	}

	// 設定を保存
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("設定の保存に失敗: %w", err)
	}

	fmt.Printf("設定 %s を %s に設定しました\n", c.Key, c.Value)
	return nil
}

// ConfigUnsetCmd は設定値を削除するコマンドです
type ConfigUnsetCmd struct {
	Key string `arg:"" help:"Configuration key to unset"`
}

// Run はconfig unsetコマンドを実行します
func (c *ConfigUnsetCmd) Run(root *RootFlags) error {
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

	// キーに応じて値を削除（空文字列に設定）
	switch c.Key {
	case "default_site":
		cfg.DefaultSite = ""
	case "keyring_backend":
		cfg.KeyringBackend = ""
	default:
		return fmt.Errorf("不明な設定キー: %s", c.Key)
	}

	// 設定を保存
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("設定の保存に失敗: %w", err)
	}

	fmt.Printf("設定 %s を削除しました\n", c.Key)
	return nil
}

// ConfigListCmd はすべての設定を一覧表示するコマンドです
type ConfigListCmd struct{}

// Run はconfig listコマンドを実行します
func (c *ConfigListCmd) Run(root *RootFlags) error {
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

	// 設定を表示
	fmt.Printf("default_site=%s\n", cfg.DefaultSite)
	fmt.Printf("keyring_backend=%s\n", cfg.KeyringBackend)

	return nil
}

// ConfigPathCmd は設定ファイルパスを表示するコマンドです
type ConfigPathCmd struct{}

// Run はconfig pathコマンドを実行します
func (c *ConfigPathCmd) Run(root *RootFlags) error {
	// 設定ファイルパスを取得
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	fmt.Println(configPath)
	return nil
}

// ConfigKeysCmd は利用可能な設定キー一覧を表示するコマンドです
type ConfigKeysCmd struct{}

// Run はconfig keysコマンドを実行します
func (c *ConfigKeysCmd) Run(root *RootFlags) error {
	fmt.Println("default_site")
	fmt.Println("keyring_backend")
	return nil
}
