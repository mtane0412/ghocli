/**
 * store_test.go
 * キーリング統合のテストコード
 */

package secrets

import (
	"os"
	"testing"
)

// TestStore_SetとGetの基本動作
func TestStore_SetとGetの基本動作(t *testing.T) {
	// テスト用のストアを作成（fileバックエンドを使用）
	store, err := NewStore("file", t.TempDir())
	if err != nil {
		t.Fatalf("ストアの作成に失敗: %v", err)
	}

	// APIキーを保存
	testKey := "64fac5417c4c6b0001234567:89abcdef01234567890123456789abcd01234567890123456789abcdef0123"
	if err := store.Set("testsite", testKey); err != nil {
		t.Fatalf("APIキーの保存に失敗: %v", err)
	}

	// APIキーを取得
	retrieved, err := store.Get("testsite")
	if err != nil {
		t.Fatalf("APIキーの取得に失敗: %v", err)
	}

	// 保存したキーと取得したキーが一致することを確認
	if retrieved != testKey {
		t.Errorf("取得したキー = %q; want %q", retrieved, testKey)
	}
}

// TestStore_存在しないキーの取得
func TestStore_存在しないキーの取得(t *testing.T) {
	store, err := NewStore("file", t.TempDir())
	if err != nil {
		t.Fatalf("ストアの作成に失敗: %v", err)
	}

	// 存在しないキーを取得しようとする
	_, err = store.Get("nonexistent")
	if err == nil {
		t.Error("存在しないキーの取得でエラーが返されなかった")
	}
}

// TestStore_Deleteでキーを削除
func TestStore_Deleteでキーを削除(t *testing.T) {
	store, err := NewStore("file", t.TempDir())
	if err != nil {
		t.Fatalf("ストアの作成に失敗: %v", err)
	}

	// APIキーを保存
	testKey := "64fac5417c4c6b0001234567:89abcdef01234567890123456789abcd01234567890123456789abcdef0123"
	if err := store.Set("testsite", testKey); err != nil {
		t.Fatalf("APIキーの保存に失敗: %v", err)
	}

	// キーを削除
	if err := store.Delete("testsite"); err != nil {
		t.Fatalf("APIキーの削除に失敗: %v", err)
	}

	// 削除後は取得できないことを確認
	_, err = store.Get("testsite")
	if err == nil {
		t.Error("削除したキーが取得できてしまった")
	}
}

// TestStore_Listで保存済みキーを一覧取得
func TestStore_Listで保存済みキーを一覧取得(t *testing.T) {
	store, err := NewStore("file", t.TempDir())
	if err != nil {
		t.Fatalf("ストアの作成に失敗: %v", err)
	}

	// 複数のAPIキーを保存
	keys := map[string]string{
		"site1": "key1:secret1",
		"site2": "key2:secret2",
		"site3": "key3:secret3",
	}
	for alias, key := range keys {
		if err := store.Set(alias, key); err != nil {
			t.Fatalf("APIキーの保存に失敗 (%s): %v", alias, err)
		}
	}

	// 保存済みのキー一覧を取得
	aliases, err := store.List()
	if err != nil {
		t.Fatalf("キー一覧の取得に失敗: %v", err)
	}

	// すべてのエイリアスが含まれていることを確認
	if len(aliases) != len(keys) {
		t.Errorf("キー数 = %d; want %d", len(aliases), len(keys))
	}

	for alias := range keys {
		found := false
		for _, a := range aliases {
			if a == alias {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("エイリアス %q が一覧に含まれていない", alias)
		}
	}
}

// TestStore_ParseAdminAPIKeyでキーをパース
func TestStore_ParseAdminAPIKeyでキーをパース(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		wantID    string
		wantSecret string
		wantErr   bool
	}{
		{
			name:       "正しいフォーマット",
			input:      "64fac5417c4c6b0001234567:89abcdef01234567890123456789abcd01234567890123456789abcdef0123",
			wantID:     "64fac5417c4c6b0001234567",
			wantSecret: "89abcdef01234567890123456789abcd01234567890123456789abcdef0123",
			wantErr:    false,
		},
		{
			name:    "コロンなし",
			input:   "64fac5417c4c6b000123456789abcdef",
			wantErr: true,
		},
		{
			name:    "空文字列",
			input:   "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, secret, err := ParseAdminAPIKey(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが返されるべきだが、nilが返された")
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if id != tc.wantID {
				t.Errorf("id = %q; want %q", id, tc.wantID)
			}
			if secret != tc.wantSecret {
				t.Errorf("secret = %q; want %q", secret, tc.wantSecret)
			}
		})
	}
}

// TestNewStore_GHO_KEYRING_BACKEND はGHO_KEYRING_BACKEND環境変数がバックエンドを上書きすることをテストします
func TestNewStore_GHO_KEYRING_BACKEND(t *testing.T) {
	// 環境変数を設定
	os.Setenv("GHO_KEYRING_BACKEND", "file")
	defer os.Unsetenv("GHO_KEYRING_BACKEND")

	// backend引数に"auto"を渡しても、環境変数でfileが使用されるはず
	store, err := NewStore("auto", t.TempDir())
	if err != nil {
		t.Fatalf("ストアの作成に失敗: %v", err)
	}

	// ストアが正しく作成されていることを確認（基本的な操作が可能）
	testKey := "test:key"
	if err := store.Set("test", testKey); err != nil {
		t.Fatalf("APIキーの保存に失敗: %v", err)
	}

	retrieved, err := store.Get("test")
	if err != nil {
		t.Fatalf("APIキーの取得に失敗: %v", err)
	}

	if retrieved != testKey {
		t.Errorf("取得したキー = %q; want %q", retrieved, testKey)
	}
}

// TestNewStore_GHO_KEYRING_PASSWORD はGHO_KEYRING_PASSWORD環境変数がパスワードを提供することをテストします
func TestNewStore_GHO_KEYRING_PASSWORD(t *testing.T) {
	// パスワード環境変数を設定
	testPassword := "test-password"
	os.Setenv("GHO_KEYRING_PASSWORD", testPassword)
	defer os.Unsetenv("GHO_KEYRING_PASSWORD")

	// fileバックエンドを使用してストアを作成
	store, err := NewStore("file", t.TempDir())
	if err != nil {
		t.Fatalf("ストアの作成に失敗: %v", err)
	}

	// ストアが正しく作成されていることを確認（基本的な操作が可能）
	testKey := "test:key"
	if err := store.Set("test", testKey); err != nil {
		t.Fatalf("APIキーの保存に失敗: %v", err)
	}

	retrieved, err := store.Get("test")
	if err != nil {
		t.Fatalf("APIキーの取得に失敗: %v", err)
	}

	if retrieved != testKey {
		t.Errorf("取得したキー = %q; want %q", retrieved, testKey)
	}
}
