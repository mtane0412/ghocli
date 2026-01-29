# gho 実装計画

## 概要

gog-cliの使用感を備えたGhost Admin APIのCLIツールを作成します。

## 技術スタック

- **言語**: Go 1.22+
- **CLIフレームワーク**: Kong (`github.com/alecthomas/kong`)
- **認証情報管理**: 99designs/keyring（OSキーリング統合）
- **JWT**: golang-jwt/jwt/v5

## 実装フェーズ

### Phase 1: 基盤構築 ✅

**目標**: プロジェクトの基盤を構築し、認証とサイト情報取得を実装する

**実装内容**:

1. **プロジェクト初期化**
   - `go mod init github.com/mtane0412/gho`
   - 依存関係追加

2. **設定システム** (`internal/config/`)
   - 設定ファイル: `~/.config/gho/config.json`
   - マルチサイト対応（エイリアス機能）
   - 構造:
     ```json
     {
       "keyring_backend": "auto",
       "default_site": "myblog",
       "sites": {
         "myblog": "https://myblog.ghost.io",
         "company": "https://blog.company.com"
       }
     }
     ```

3. **キーリング統合** (`internal/secrets/`)
   - Admin APIキーの安全な保存
   - バックエンド: auto/file/keychain/secretservice/wincred

4. **Ghost APIクライアント** (`internal/ghostapi/`)
   - JWT生成（HS256、有効期限5分）
   - HTTPクライアント
   - サイト情報取得API

5. **出力フォーマット** (`internal/outfmt/`)
   - JSON形式
   - テーブル形式
   - TSV形式（プレーン）

6. **認証コマンド** (`internal/cmd/auth.go`)
   - `gho auth add <site-url>` - APIキー登録
   - `gho auth list` - 登録済みサイト一覧
   - `gho auth remove <alias>` - APIキー削除
   - `gho auth status` - 認証状態確認

7. **基本コマンド**
   - `gho site` - サイト情報取得
   - `gho version` - バージョン表示

**検証方法**:
```bash
make build
./gho auth add https://your-ghost-site.ghost.io
./gho auth list
./gho site
```

**完了**: ✅ 2026-01-29

---

### Phase 2: コンテンツ管理（Posts/Pages）

**目標**: Posts/Pagesの作成、更新、削除、公開機能を実装する

**実装内容**:

1. **Posts API** (`internal/ghostapi/posts.go`)
   - `ListPosts(options ListOptions) ([]Post, error)`
   - `GetPost(idOrSlug string) (*Post, error)`
   - `CreatePost(post *Post) (*Post, error)`
   - `UpdatePost(id string, post *Post) (*Post, error)`
   - `DeletePost(id string) error`

2. **Posts型定義**
   ```go
   type Post struct {
       ID          string     `json:"id"`
       Title       string     `json:"title"`
       Slug        string     `json:"slug"`
       HTML        string     `json:"html,omitempty"`
       MobileDoc   string     `json:"mobiledoc,omitempty"`
       Status      string     `json:"status"` // draft/published/scheduled
       CreatedAt   time.Time  `json:"created_at"`
       PublishedAt *time.Time `json:"published_at,omitempty"`
       Tags        []Tag      `json:"tags,omitempty"`
       Authors     []Author   `json:"authors,omitempty"`
   }
   ```

3. **Postsコマンド** (`internal/cmd/posts.go`)
   - `gho posts list [--status draft|published|scheduled] [--limit N]`
   - `gho posts get <id-or-slug>`
   - `gho posts create --title "..." [--html "..."]`
   - `gho posts update <id> --title "..."`
   - `gho posts delete <id>`
   - `gho posts publish <id>`

4. **Pages API** (`internal/ghostapi/pages.go`)
   - `ListPages(options ListOptions) ([]Page, error)`
   - `GetPage(idOrSlug string) (*Page, error)`
   - `CreatePage(page *Page) (*Page, error)`
   - `UpdatePage(id string, page *Page) (*Page, error)`
   - `DeletePage(id string) error`

5. **Pagesコマンド** (`internal/cmd/pages.go`)
   - `gho pages list`
   - `gho pages get <id-or-slug>`
   - `gho pages create --title "..."`
   - `gho pages update <id> ...`
   - `gho pages delete <id>`

**テスト**:
- Posts APIのテスト（`internal/ghostapi/posts_test.go`）
- Postsコマンドのテスト（`internal/cmd/posts_test.go`）
- Pages APIのテスト（`internal/ghostapi/pages_test.go`）
- Pagesコマンドのテスト（`internal/cmd/pages_test.go`）

**検証方法**:
```bash
./gho posts list
./gho posts get <slug>
./gho posts create --title "Test Post" --status draft
./gho posts publish <id>
./gho posts delete <id>

./gho pages list
./gho pages create --title "Test Page"
```

---

### Phase 3: タクソノミー + メディア

**目標**: Tags管理とImages アップロード機能を実装する

**実装内容**:

1. **Tags API** (`internal/ghostapi/tags.go`)
   - `ListTags() ([]Tag, error)`
   - `GetTag(idOrSlug string) (*Tag, error)`
   - `CreateTag(tag *Tag) (*Tag, error)`
   - `UpdateTag(id string, tag *Tag) (*Tag, error)`
   - `DeleteTag(id string) error`

2. **Tags型定義**
   ```go
   type Tag struct {
       ID          string `json:"id"`
       Name        string `json:"name"`
       Slug        string `json:"slug"`
       Description string `json:"description,omitempty"`
   }
   ```

3. **Tagsコマンド** (`internal/cmd/tags.go`)
   - `gho tags list`
   - `gho tags get <id-or-slug>`
   - `gho tags create --name "..."`
   - `gho tags update <id> --name "..."`
   - `gho tags delete <id>`

4. **Images API** (`internal/ghostapi/images.go`)
   - `UploadImage(filePath string) (*ImageUploadResponse, error)`

5. **Imagesコマンド** (`internal/cmd/images.go`)
   - `gho images upload <file-path>`

**テスト**:
- Tags APIのテスト
- Tagsコマンドのテスト
- Images APIのテスト
- Imagesコマンドのテスト

**検証方法**:
```bash
./gho tags list
./gho tags create --name "Technology"
./gho images upload ./image.png
```

---

### Phase 4: Members管理

**目標**: Members（購読者）の管理機能を実装する

**実装内容**:

1. **Members API** (`internal/ghostapi/members.go`)
   - `ListMembers(options ListOptions) ([]Member, error)`
   - `GetMember(id string) (*Member, error)`
   - `CreateMember(member *Member) (*Member, error)`
   - `UpdateMember(id string, member *Member) (*Member, error)`
   - `DeleteMember(id string) error`

2. **Membersコマンド** (`internal/cmd/members.go`)
   - `gho members list`
   - `gho members get <id>`
   - `gho members create --email "..."`
   - `gho members update <id> ...`
   - `gho members delete <id>`

---

### Phase 5: Users管理

**目標**: Users（管理者・編集者）の管理機能を実装する

**実装内容**:

1. **Users API** (`internal/ghostapi/users.go`)
   - `ListUsers() ([]User, error)`
   - `GetUser(id string) (*User, error)`
   - `UpdateUser(id string, user *User) (*User, error)`

2. **Usersコマンド** (`internal/cmd/users.go`)
   - `gho users list`
   - `gho users get <id>`
   - `gho users update <id> ...`

---

### Phase 6: Newsletters/Tiers/Offers

**目標**: Newsletter、Tiers（購読プラン）、Offers（特典）の管理機能を実装する

**実装内容**:

1. **Newsletters API**
   - `ListNewsletters() ([]Newsletter, error)`
   - `GetNewsletter(id string) (*Newsletter, error)`

2. **Tiers API**
   - `ListTiers() ([]Tier, error)`
   - `GetTier(id string) (*Tier, error)`

3. **Offers API**
   - `ListOffers() ([]Offer, error)`
   - `GetOffer(id string) (*Offer, error)`

---

### Phase 7: Themes/Webhooks

**目標**: Themes管理とWebhooks管理機能を実装する

**実装内容**:

1. **Themes API**
   - `ListThemes() ([]Theme, error)`
   - `UploadTheme(filePath string) error`
   - `ActivateTheme(name string) error`
   - `DeleteTheme(name string) error`

2. **Webhooks API**
   - `ListWebhooks() ([]Webhook, error)`
   - `CreateWebhook(webhook *Webhook) (*Webhook, error)`
   - `DeleteWebhook(id string) error`

---

## 開発ワークフロー

### TDD原則

すべての実装において、以下のTDDサイクルに従います：

1. **RED** - 失敗するテストを先に書く
2. **GREEN** - テストを通す最小限のコードを書く
3. **REFACTOR** - コードを整理する

### 品質チェック

各フェーズ完了時に以下を実行：

```bash
# テスト実行
make test

# 型チェック
make type-check

# Lint実行
make lint

# ビルド確認
make build
```

### Git ワークフロー

```bash
# フェーズごとに feature ブランチを作成
git checkout -b feature/phase2-content-management

# コミット前チェック
make test
make type-check
make lint

# コミット作成
git add .
git commit -m "Phase 2: コンテンツ管理機能を実装

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## 参照リソース

### Ghost Admin API ドキュメント
- https://ghost.org/docs/admin-api/

### 参照プロジェクト（gog-cli）
- `../gogcli/internal/cmd/root.go` - CLI構造体パターン
- `../gogcli/internal/cmd/auth.go` - 認証コマンド実装
- `../gogcli/internal/secrets/store.go` - キーリング統合
- `../gogcli/internal/config/config.go` - 設定ファイル管理
- `../gogcli/internal/outfmt/outfmt.go` - 出力フォーマット
