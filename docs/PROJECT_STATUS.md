# プロジェクト状態

## 概要

**gho** はGhost Admin APIのCLIツールです。gog-cliの使用感を備え、Ghost Admin APIの操作をコマンドラインから実行できます。

## 実装フェーズ

### ✅ Phase 1: 基盤構築（完了）

**完了日**: 2026-01-29

**実装内容**:

1. **プロジェクト初期化**
   - Go modules初期化
   - 依存関係追加（Kong、Keyring、JWT）

2. **設定システム** (`internal/config/`)
   - 設定ファイル管理（`~/.config/gho/config.json`）
   - マルチサイト対応（エイリアス機能）
   - デフォルトサイト管理

3. **キーリング統合** (`internal/secrets/`)
   - OSキーリングによる安全なAPIキー保存
   - macOS Keychain、Linux Secret Service、Windows Credential Manager対応
   - APIキーのパース機能

4. **Ghost APIクライアント** (`internal/ghostapi/`)
   - JWT生成機能（HS256、有効期限5分）
   - HTTPクライアント
   - サイト情報取得API

5. **出力フォーマット** (`internal/outfmt/`)
   - JSON形式
   - テーブル形式（人間向け）
   - TSV形式（プログラム連携向け）

6. **認証コマンド** (`internal/cmd/auth.go`)
   ```
   gho auth add <site-url>      # APIキー登録
   gho auth list                # 登録済みサイト一覧
   gho auth remove <alias>      # APIキー削除
   gho auth status              # 認証状態確認
   ```

7. **基本コマンド**
   ```
   gho site                     # サイト情報取得
   gho version                  # バージョン表示
   ```

**品質チェック**:
- ✅ すべてのテストがパス
- ✅ 型チェック（`go vet`）成功
- ✅ ビルド成功

**コミット**: `68b9340 Phase 1: 基盤実装を完了`

### ✅ Phase 2: コンテンツ管理（Posts/Pages）（完了）

**完了日**: 2026-01-29

**実装内容**:

1. **Posts API** (`internal/ghostapi/posts.go`)
   - Post型定義（ID、Title、Slug、HTML、Status、PublishedAtなど）
   - ListOptions型定義（Limit、Status、Filterなど）
   - `ListPosts(options ListOptions) ([]Post, error)` 実装
   - `GetPost(idOrSlug string) (*Post, error)` 実装
   - `CreatePost(post *Post) (*Post, error)` 実装
   - `UpdatePost(id string, post *Post) (*Post, error)` 実装
   - `DeletePost(id string) error` 実装

2. **Pages API** (`internal/ghostapi/pages.go`)
   - Page型定義（ID、Title、Slug、HTML、Statusなど）
   - `ListPages(options ListOptions) ([]Page, error)` 実装
   - `GetPage(idOrSlug string) (*Page, error)` 実装
   - `CreatePage(page *Page) (*Page, error)` 実装
   - `UpdatePage(id string, page *Page) (*Page, error)` 実装
   - `DeletePage(id string) error` 実装

3. **Postsコマンド** (`internal/cmd/posts.go`)
   ```
   gho posts list [--status draft|published|scheduled] [--limit N]
   gho posts get <id-or-slug>
   gho posts create --title "..." [--html "..."] [--status draft|published]
   gho posts update <id> [--title "..."] [--html "..."]
   gho posts delete <id>
   gho posts publish <id>
   ```

4. **Pagesコマンド** (`internal/cmd/pages.go`)
   ```
   gho pages list [--status draft|published|scheduled] [--limit N]
   gho pages get <id-or-slug>
   gho pages create --title "..." [--html "..."]
   gho pages update <id> [--title "..."] [--html "..."]
   gho pages delete <id>
   ```

**品質チェック**:
- ✅ すべてのテストがパス（Posts: 7テスト、Pages: 5テスト）
- ✅ 型チェック（`go vet`）成功
- ✅ ビルド成功

**コミット**:
- `40c33f2 feat(ghostapi): Posts APIを実装`
- `016fe5c feat(ghostapi): Pages APIを実装`
- `a84e3da feat(cmd): Posts/Pagesコマンドを実装`

### ✅ Phase 3: タクソノミー + メディア（完了）

**完了日**: 2026-01-30

**実装内容**:

1. **Tags API** (`internal/ghostapi/tags.go`)
   - Tag型定義（ID、Name、Slug、Description、Visibilityなど）
   - TagListOptions型定義（pagination、filter対応）
   - `ListTags(options TagListOptions) (*TagListResponse, error)` 実装
   - `GetTag(idOrSlug string) (*Tag, error)` 実装（"slug:"プレフィックス対応）
   - `CreateTag(tag *Tag) (*Tag, error)` 実装
   - `UpdateTag(id string, tag *Tag) (*Tag, error)` 実装
   - `DeleteTag(id string) error` 実装

2. **Images API** (`internal/ghostapi/images.go`)
   - Image型定義（URL、Ref）
   - `UploadImage(file io.Reader, filename string, opts ImageUploadOptions) (*Image, error)` 実装
   - multipart/form-dataでのアップロード対応
   - Purpose（image/profile_image/icon）指定対応

3. **Tagsコマンド** (`internal/cmd/tags.go`)
   ```
   gho tags list [--limit N] [--page N]
   gho tags get <id-or-slug>        # "slug:tag-name" 形式でslugを指定可能
   gho tags create --name "..." [--description "..."] [--visibility public|internal]
   gho tags update <id> [--name "..."] [--description "..."]
   gho tags delete <id>
   ```

4. **Imagesコマンド** (`internal/cmd/images.go`)
   ```
   gho images upload <file-path> [--purpose image|profile_image|icon] [--ref <ref-id>]
   ```

**品質チェック**:
- ✅ すべてのテストがパス（Tags: 6テスト、Images: 2テスト）
- ✅ 型チェック（`go vet`）成功
- ✅ ビルド成功

**コミット**:
- `b5299e8 feat(api): Tags APIとImages APIを実装`

### 📋 Phase 4以降（未実装）

- Members API
- Users API
- Newsletters API
- Tiers API
- Offers API
- Themes API
- Webhooks API

## 現在の構造

```
gho/
├── cmd/gho/
│   └── main.go              # エントリーポイント
├── internal/
│   ├── cmd/                  # CLIコマンド定義
│   │   ├── root.go          # CLI構造体、RootFlags
│   │   ├── auth.go          # 認証コマンド
│   │   ├── site.go          # サイト情報コマンド
│   │   ├── posts.go         # Postsコマンド
│   │   ├── pages.go         # Pagesコマンド
│   │   ├── tags.go          # Tagsコマンド
│   │   └── images.go        # Imagesコマンド
│   ├── config/              # 設定ファイル管理
│   │   ├── config.go
│   │   └── config_test.go
│   ├── secrets/             # キーリング統合
│   │   ├── store.go
│   │   └── store_test.go
│   ├── ghostapi/            # Ghost APIクライアント
│   │   ├── client.go        # HTTPクライアント
│   │   ├── client_test.go
│   │   ├── jwt.go           # JWT生成
│   │   ├── jwt_test.go
│   │   ├── posts.go         # Posts API
│   │   ├── posts_test.go
│   │   ├── pages.go         # Pages API
│   │   ├── pages_test.go
│   │   ├── tags.go          # Tags API
│   │   ├── tags_test.go
│   │   ├── images.go        # Images API
│   │   └── images_test.go
│   └── outfmt/              # 出力フォーマット
│       ├── outfmt.go
│       └── outfmt_test.go
├── docs/                    # ドキュメント
├── go.mod
├── go.sum
├── Makefile
├── .golangci.yml
├── .gitignore
└── README.md
```

## テストカバレッジ

すべてのコアコンポーネントはテスト済みです：

- `internal/config/` - 設定ファイル管理（6テスト）
- `internal/secrets/` - キーリング統合（8テスト）
- `internal/ghostapi/` - APIクライアント（29テスト）
  - `client.go`, `jwt.go` - 9テスト
  - `posts.go` - 7テスト
  - `pages.go` - 5テスト
  - `tags.go` - 6テスト
  - `images.go` - 2テスト
- `internal/outfmt/` - 出力フォーマット（5テスト）

合計: 48テスト、すべてパス

## 依存関係

```
github.com/alecthomas/kong v1.13.0        # CLIフレームワーク
github.com/99designs/keyring v1.2.2       # キーリング統合
github.com/golang-jwt/jwt/v5 v5.3.1       # JWT生成
```

## 品質チェックコマンド

```bash
# テスト実行
make test

# 型チェック
make type-check

# Lint実行（golangci-lintが必要）
make lint

# ビルド
make build
```

## 次のステップ

Phase 4（Members管理）の実装を開始します。詳細は `docs/NEXT_STEPS.md` を参照してください。
