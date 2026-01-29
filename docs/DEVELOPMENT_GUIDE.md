# gho 開発ガイド

## 開発環境のセットアップ

### 必要なツール

- **Go**: 1.22以上
- **Make**: ビルド自動化
- **golangci-lint**: Lint実行（オプション）
- **Git**: バージョン管理

### 環境構築手順

```bash
# リポジトリをクローン
git clone https://github.com/mtane0412/gho.git
cd gho

# 依存関係をインストール
go mod download

# ビルド
make build

# テスト実行
make test
```

## 開発ワークフロー

### TDD原則

すべての実装において、以下のTDDサイクルに従います：

1. **RED** - 失敗するテストを先に書く
2. **GREEN** - テストを通す最小限のコードを書く
3. **REFACTOR** - コードを整理する

### 実装例

```go
// 1. RED: 失敗するテストを先に書く
func TestGenerateJWT_正しいフォーマットのトークン生成(t *testing.T) {
    token, err := GenerateJWT("keyid", "secret")
    if err != nil {
        t.Fatalf("JWTの生成に失敗: %v", err)
    }
    if token == "" {
        t.Error("生成されたトークンが空です")
    }
}

// 2. GREEN: テストを通す最小限のコードを書く
func GenerateJWT(keyID, secret string) (string, error) {
    // 最小限の実装
}

// 3. REFACTOR: コードを整理する
func GenerateJWT(keyID, secret string) (string, error) {
    // リファクタリング後の実装
}
```

## コーディング規約

### ファイル冒頭コメント

各ファイルの冒頭に仕様をコメントで記述する：

```go
/**
 * jwt.go
 * Ghost Admin API用のJWT生成
 *
 * Ghost Admin APIはHS256アルゴリズムで署名されたJWTを要求します。
 * トークンの有効期限は5分です。
 */

package ghostapi
```

### 関数コメント

目的、内容、注意事項を詳細に日本語で記述する：

```go
// GenerateJWT はGhost Admin API用のJWTトークンを生成します。
// keyID: Admin APIキーのID部分
// secret: Admin APIキーのシークレット部分
func GenerateJWT(keyID, secret string) (string, error) {
    // ...
}
```

### テスト関数の命名

テスト関数名は日本語で具体的な内容を記述：

```go
// ✅ 良い例
func TestGenerateJWT_正しいフォーマットのトークン生成(t *testing.T) { }
func TestGenerateJWT_空のキーIDでエラー(t *testing.T) { }

// ❌ 悪い例
func TestGenerateJWT(t *testing.T) { }
func TestJWT1(t *testing.T) { }
```

### エラーハンドリング

エラーメッセージは日本語で具体的に：

```go
// ✅ 良い例
if keyID == "" {
    return "", errors.New("キーIDが空です")
}

// ❌ 悪い例
if keyID == "" {
    return "", errors.New("invalid key")
}
```

エラーのラップには `fmt.Errorf` と `%w` を使用：

```go
if err := store.Set(alias, apiKey); err != nil {
    return fmt.Errorf("APIキーの保存に失敗: %w", err)
}
```

### 構造体タグ

JSON構造体には適切なタグを付与：

```go
type Site struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    URL         string `json:"url"`
    Version     string `json:"version"`
}
```

## テスト

### テスト実行

```bash
# すべてのテストを実行
go test ./...

# 詳細表示
go test ./... -v

# 特定のパッケージのみ
go test ./internal/config/... -v

# カバレッジ付き
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### テストの書き方

#### ユニットテスト

```go
func TestStore_SetとGetの基本動作(t *testing.T) {
    // テスト用のストアを作成（fileバックエンドを使用）
    store, err := NewStore("file", t.TempDir())
    if err != nil {
        t.Fatalf("ストアの作成に失敗: %v", err)
    }

    // APIキーを保存
    testKey := "64fac5417c4c6b0001234567:89abcdef..."
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
```

#### HTTPクライアントのテスト

`httptest` パッケージを使用してモックサーバーを作成：

```go
func TestGetSite_サイト情報の取得(t *testing.T) {
    // テスト用のHTTPサーバーを作成
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // レスポンスを返す
        response := map[string]interface{}{
            "site": map[string]interface{}{
                "title": "Test Blog",
                "url":   "https://test.ghost.io",
            },
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }))
    defer server.Close()

    // クライアントを作成
    client, err := NewClient(server.URL, "keyid", "secret")
    if err != nil {
        t.Fatalf("クライアントの作成に失敗: %v", err)
    }

    // サイト情報を取得
    site, err := client.GetSite()
    if err != nil {
        t.Fatalf("サイト情報の取得に失敗: %v", err)
    }

    // レスポンスの検証
    if site.Title != "Test Blog" {
        t.Errorf("Title = %q; want %q", site.Title, "Test Blog")
    }
}
```

#### テーブル駆動テスト

複数のテストケースを効率的にテスト：

```go
func TestParseAdminAPIKey(t *testing.T) {
    testCases := []struct {
        name       string
        input      string
        wantID     string
        wantSecret string
        wantErr    bool
    }{
        {
            name:       "正しいフォーマット",
            input:      "64fac5417c4c6b0001234567:89abcdef...",
            wantID:     "64fac5417c4c6b0001234567",
            wantSecret: "89abcdef...",
            wantErr:    false,
        },
        {
            name:    "コロンなし",
            input:   "64fac5417c4c6b000123456789abcdef",
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
```

## 品質チェック

### コミット前チェック

コミット前に以下を必ず実行：

```bash
# テスト実行
make test

# 型チェック
make type-check

# Lint実行（golangci-lintが必要）
make lint

# ビルド確認
make build
```

### 型チェック

```bash
# go vet で型チェック
go vet ./...

# または
make type-check
```

### Lint

golangci-lintをインストール：

```bash
# macOS
brew install golangci-lint

# Linux/Windows
# https://golangci-lint.run/usage/install/
```

Lint実行：

```bash
golangci-lint run

# または
make lint
```

## Git ワークフロー

### ブランチ戦略

```bash
# mainブランチへの直接コミット禁止
# 必ずfeatureブランチで作業する

# 作業開始前にブランチを確認
git branch --show-current

# featureブランチを作成
git checkout -b feature/phase2-content-management
```

### コミットメッセージ

```bash
git commit -m "$(cat <<'EOF'
Phase 2: コンテンツ管理機能を実装

Posts/Pagesの作成、更新、削除、公開機能を実装しました。

主な実装内容：
- Posts API（list/get/create/update/delete/publish）
- Pages API（list/get/create/update/delete）
- Posts/Pagesコマンド
- テストの追加

すべてのテストがパスしています。

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
EOF
)"
```

## 新しいAPIリソースの追加方法

### 1. API型定義を追加

`internal/ghostapi/posts.go`:

```go
package ghostapi

import "time"

// Post はGhostの投稿を表します
type Post struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Slug        string     `json:"slug"`
    HTML        string     `json:"html,omitempty"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    PublishedAt *time.Time `json:"published_at,omitempty"`
}
```

### 2. テストを先に書く（RED）

`internal/ghostapi/posts_test.go`:

```go
func TestListPosts_投稿一覧の取得(t *testing.T) {
    // テスト用のHTTPサーバーを作成
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // モックレスポンスを返す
    }))
    defer server.Close()

    // テストを書く
}
```

### 3. 実装を追加（GREEN）

`internal/ghostapi/posts.go`:

```go
// ListPosts は投稿一覧を取得します
func (c *Client) ListPosts(options ListOptions) ([]Post, error) {
    // 実装
}
```

### 4. コマンドを追加

`internal/cmd/posts.go`:

```go
type PostsCmd struct {
    List   PostsListCmd   `cmd:"" help:"List posts"`
    Get    PostsGetCmd    `cmd:"" help:"Get a post"`
    Create PostsCreateCmd `cmd:"" help:"Create a post"`
}
```

### 5. root.goに登録

`internal/cmd/root.go`:

```go
type CLI struct {
    RootFlags `embed:""`
    Version   kong.VersionFlag
    Auth      AuthCmd
    Site      SiteCmd
    Posts     PostsCmd  // 追加
}
```

## デバッグ

### ログ出力

`--verbose` フラグでログを有効化：

```go
if root.Verbose {
    log.Printf("APIリクエスト: %s %s", method, url)
}
```

### JWTのデバッグ

jwt.ioでトークンをデコード：

```bash
# トークンを取得
./gho site --verbose

# jwt.io にアクセスして貼り付け
```

### HTTPリクエストのデバッグ

環境変数で詳細ログを有効化：

```bash
# HTTPデバッグ
export GODEBUG=http2debug=1
./gho site
```

## トラブルシューティング

### テストが失敗する

```bash
# キャッシュをクリア
go clean -testcache

# 再実行
go test ./...
```

### ビルドエラー

```bash
# 依存関係を更新
go mod tidy

# 再ビルド
make build
```

### キーリングエラー

```bash
# fileバックエンドで動作確認
export GHO_KEYRING_BACKEND=file
./gho auth add https://test.ghost.io
```

## リリース

### バージョンタグ

```bash
# タグを作成
git tag -a v0.1.0 -m "Release v0.1.0"

# タグをプッシュ
git push origin v0.1.0
```

### ビルド

```bash
# バージョン指定ビルド
make build VERSION=0.1.0
```

## 参考リソース

### Go言語

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Ghost Admin API

- [Ghost Admin API Documentation](https://ghost.org/docs/admin-api/)
- [Ghost API Client Examples](https://github.com/TryGhost/Ghost/tree/main/ghost/admin-api)

### テスト

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
