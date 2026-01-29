/**
 * store.go
 * キーリング統合によるAdmin APIキーの安全な保存
 *
 * 99designs/keyringを使用してOSのキーリングにAPIキーを保存します。
 * macOS: Keychain、Linux: Secret Service、Windows: Credential Managerに対応。
 */

package secrets

import (
	"errors"
	"fmt"
	"strings"

	"github.com/99designs/keyring"
)

// Store はAPIキーを保存・取得するためのストアです
type Store struct {
	ring keyring.Keyring
}

const (
	// ServiceName はキーリングに保存する際のサービス名
	ServiceName = "gho-ghost-admin"
)

// NewStore は新しいキーリングストアを作成します。
// backend: "auto", "file", "keychain" など
// fileDir: backendが"file"の場合のファイル保存先ディレクトリ
func NewStore(backend, fileDir string) (*Store, error) {
	var cfg keyring.Config
	cfg.ServiceName = ServiceName

	// バックエンドタイプを設定
	switch backend {
	case "auto":
		cfg.AllowedBackends = []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.SecretServiceBackend,
			keyring.WinCredBackend,
			keyring.FileBackend,
		}
	case "file":
		cfg.AllowedBackends = []keyring.BackendType{keyring.FileBackend}
		cfg.FileDir = fileDir
		cfg.FilePasswordFunc = func(prompt string) (string, error) {
			// ファイルバックエンドはパスワードなしで使用
			return "", nil
		}
	case "keychain":
		cfg.AllowedBackends = []keyring.BackendType{keyring.KeychainBackend}
	case "secretservice":
		cfg.AllowedBackends = []keyring.BackendType{keyring.SecretServiceBackend}
	case "wincred":
		cfg.AllowedBackends = []keyring.BackendType{keyring.WinCredBackend}
	default:
		return nil, fmt.Errorf("不明なバックエンド: %s", backend)
	}

	ring, err := keyring.Open(cfg)
	if err != nil {
		return nil, fmt.Errorf("キーリングのオープンに失敗: %w", err)
	}

	return &Store{ring: ring}, nil
}

// Set はサイトエイリアスに対してAdmin APIキーを保存します。
func (s *Store) Set(alias, apiKey string) error {
	item := keyring.Item{
		Key:  alias,
		Data: []byte(apiKey),
	}

	if err := s.ring.Set(item); err != nil {
		return fmt.Errorf("APIキーの保存に失敗: %w", err)
	}

	return nil
}

// Get はサイトエイリアスに対応するAdmin APIキーを取得します。
func (s *Store) Get(alias string) (string, error) {
	item, err := s.ring.Get(alias)
	if err != nil {
		return "", fmt.Errorf("APIキーの取得に失敗: %w", err)
	}

	return string(item.Data), nil
}

// Delete はサイトエイリアスに対応するAdmin APIキーを削除します。
func (s *Store) Delete(alias string) error {
	if err := s.ring.Remove(alias); err != nil {
		return fmt.Errorf("APIキーの削除に失敗: %w", err)
	}

	return nil
}

// List はキーリングに保存されているすべてのサイトエイリアスを一覧取得します。
func (s *Store) List() ([]string, error) {
	keys, err := s.ring.Keys()
	if err != nil {
		return nil, fmt.Errorf("キー一覧の取得に失敗: %w", err)
	}

	return keys, nil
}

// ParseAdminAPIKey はGhost Admin APIキー（id:secret形式）をパースします。
func ParseAdminAPIKey(apiKey string) (id, secret string, err error) {
	if apiKey == "" {
		return "", "", errors.New("APIキーが空です")
	}

	parts := strings.SplitN(apiKey, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("APIキーのフォーマットが不正です（id:secret形式である必要があります）")
	}

	return parts[0], parts[1], nil
}
