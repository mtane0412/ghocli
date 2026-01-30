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

### ✅ Phase 4: Members管理（完了）

**完了日**: 2026-01-30

**実装内容**:

1. **Members API** (`internal/ghostapi/members.go`)
   - Member型定義（ID、UUID、Email、Name、Note、Status、Labelsなど）
   - Label型定義（ID、Name、Slug）
   - MemberListOptions型定義（pagination、filter、order対応）
   - `ListMembers(options MemberListOptions) (*MemberListResponse, error)` 実装
   - `GetMember(id string) (*Member, error)` 実装
   - `CreateMember(member *Member) (*Member, error)` 実装
   - `UpdateMember(id string, member *Member) (*Member, error)` 実装
   - `DeleteMember(id string) error` 実装

2. **Membersコマンド** (`internal/cmd/members.go`)
   ```
   gho members list [--limit N] [--page N] [--filter "..."] [--order "..."]
   gho members get <id>
   gho members create --email "..." [--name "..."] [--note "..."] [--labels "..."]
   gho members update <id> [--name "..."] [--note "..."] [--labels "..."]
   gho members delete <id>
   ```

**品質チェック**:
- ✅ すべてのテストがパス（Members: 6テスト）
- ✅ 型チェック（`go vet`）成功
- ✅ ビルド成功

**コミット**:
- `3a935e6 feat(api): Members APIを実装`

### ✅ Phase 5: Users管理（完了）

**完了日**: 2026-01-30

**実装内容**:

1. **Users API** (`internal/ghostapi/users.go`)
   - User型定義（ID、Name、Slug、Email、Bio、Location、Website、ProfileImage、CoverImage、Rolesなど）
   - Role型定義（ID、Name）
   - UserListOptions型定義（pagination、include、filter対応）
   - `ListUsers(options UserListOptions) (*UserListResponse, error)` 実装
   - `GetUser(idOrSlug string) (*User, error)` 実装（"slug:"プレフィックス対応）
   - `UpdateUser(id string, user *User) (*User, error)` 実装
   - **注意**: Create/Delete操作は非サポート（Ghostダッシュボードの招待機能を利用）

2. **Usersコマンド** (`internal/cmd/users.go`)
   ```
   gho users list [--limit N] [--page N] [--include roles,count.posts]
   gho users get <id-or-slug>       # "slug:user-slug" 形式でslugを指定可能
   gho users update <id> [--name "..."] [--slug "..."] [--bio "..."] [--location "..."] [--website "..."]
   ```

**品質チェック**:
- ✅ すべてのテストがパス（Users: 7テスト）
- ✅ 型チェック（`go vet`）成功
- ✅ ビルド成功

**コミット**:
- `1884ff0 feat(api): Users APIを実装`

### ✅ Phase 6: Newsletters/Tiers/Offers（完了）

**完了日**: 2026-01-30

**実装内容**:

1. **Newsletters API** (`internal/ghostapi/newsletters.go`)
   - Newsletter型定義（ID、Name、Slug、Description、Status、SubscribeOnSignupなど）
   - NewsletterListOptions型定義（pagination、filter対応）
   - `ListNewsletters(options NewsletterListOptions) (*NewsletterListResponse, error)` 実装
   - `GetNewsletter(idOrSlug string) (*Newsletter, error)` 実装（"slug:"プレフィックス対応）
   - `CreateNewsletter(newsletter *Newsletter) (*Newsletter, error)` 実装
   - `UpdateNewsletter(id string, newsletter *Newsletter) (*Newsletter, error)` 実装

2. **Tiers API** (`internal/ghostapi/tiers.go`)
   - Tier型定義（ID、Name、Slug、Type、MonthlyPrice、YearlyPriceなど）
   - TierListOptions型定義（pagination、include対応）
   - `ListTiers(options TierListOptions) (*TierListResponse, error)` 実装
   - `GetTier(idOrSlug string) (*Tier, error)` 実装（"slug:"プレフィックス対応）
   - `CreateTier(tier *Tier) (*Tier, error)` 実装
   - `UpdateTier(id string, tier *Tier) (*Tier, error)` 実装

3. **Offers API** (`internal/ghostapi/offers.go`)
   - Offer型定義（ID、Name、Code、Tier、DiscountType、DiscountAmountなど）
   - OfferListOptions型定義（pagination、filter対応）
   - `ListOffers(options OfferListOptions) (*OfferListResponse, error)` 実装
   - `GetOffer(id string) (*Offer, error)` 実装
   - `CreateOffer(offer *Offer) (*Offer, error)` 実装
   - `UpdateOffer(id string, offer *Offer) (*Offer, error)` 実装

4. **Newslettersコマンド** (`internal/cmd/newsletters.go`)
   ```
   gho newsletters list [--limit N] [--page N] [--filter "..."]
   gho newsletters get <id-or-slug>    # "slug:newsletter-slug" 形式でslugを指定可能
   gho newsletters create --name "..." [--description "..."] [--visibility members|paid]
   gho newsletters update <id> [--name "..."] [--visibility "..."] [--sender-name "..."]
   ```

5. **Tiersコマンド** (`internal/cmd/tiers.go`)
   ```
   gho tiers list [--limit N] [--page N] [--include monthly_price,yearly_price]
   gho tiers get <id-or-slug>          # "slug:tier-slug" 形式でslugを指定可能
   gho tiers create --name "..." [--type free|paid] [--monthly-price N] [--yearly-price N]
   gho tiers update <id> [--name "..."] [--monthly-price N] [--yearly-price N]
   ```

6. **Offersコマンド** (`internal/cmd/offers.go`)
   ```
   gho offers list [--limit N] [--page N] [--filter "..."]
   gho offers get <id>
   gho offers create --name "..." --code "..." --type percent|fixed --amount N --tier-id <tier-id>
   gho offers update <id> [--name "..."] [--amount N]
   ```

7. **破壊的操作の確認機構** (`internal/cmd/helpers.go`)
   - Create/Update操作には確認プロンプトが表示される
   - `--force`フラグで確認をスキップ可能

**品質チェック**:
- ✅ すべてのテストがパス（Newsletters: 6テスト、Tiers: 6テスト、Offers: 6テスト）
- ✅ 型チェック（`go vet`）成功
- ✅ ビルド成功

**コミット**:
- `4545035 feat(api): Newsletters, Tiers, Offers APIを実装`
- `eed5ff2 feat(newsletters): Newsletters書き込み操作を実装`
- `8b158df feat(tiers): Tiers書き込み操作を実装`
- `2874e8d feat(offers): Offers書き込み操作を実装`
- `013086c feat(cmd): 破壊的操作の確認機構を実装`

### ✅ Phase 7: Themes/Webhooks API（完了）

**完了日**: 2026-01-30

**実装内容**:

1. **Themes API** (`internal/ghostapi/themes.go`)
   - Theme型定義（Name、Package、Active、Templatesなど）
   - ThemePackage型定義（Name、Description、Version）
   - ThemeTemplate型定義（Filename）
   - `ListThemes() (*ThemeListResponse, error)` 実装
   - `UploadTheme(file io.Reader, filename string) (*Theme, error)` 実装（multipartアップロード）
   - `ActivateTheme(name string) (*Theme, error)` 実装

2. **Webhooks API** (`internal/ghostapi/webhooks.go`)
   - Webhook型定義（ID、Event、TargetURL、Name、Secret、APIVersion、IntegrationID、Status、LastTriggeredAt、CreatedAt、UpdatedAtなど）
   - `CreateWebhook(webhook *Webhook) (*Webhook, error)` 実装
   - `UpdateWebhook(id string, webhook *Webhook) (*Webhook, error)` 実装
   - `DeleteWebhook(id string) error` 実装
   - **注意**: Ghost APIはWebhookのList/Getをサポートしていません

3. **Themesコマンド** (`internal/cmd/themes.go`)
   ```
   gho themes list                    # テーマ一覧
   gho themes upload <file.zip>       # テーマアップロード
   gho themes activate <name>         # テーマ有効化
   ```

4. **Webhooksコマンド** (`internal/cmd/webhooks.go`)
   ```
   gho webhooks create --event <event> --target-url <url> [--name <name>]
   gho webhooks update <id> [--event <event>] [--target-url <url>] [--name <name>]
   gho webhooks delete <id>
   ```

**品質チェック**:
- ✅ すべてのテストがパス（Themes: 3テスト、Webhooks: 3テスト）
- ✅ 型チェック（`go vet`）成功
- ✅ Lint（golangci-lint）成功
- ✅ ビルド成功

**コミット**:
- `3a6a9ed feat(api): Themes/Webhooks APIを実装`

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
│   │   ├── images.go        # Imagesコマンド
│   │   ├── members.go       # Membersコマンド
│   │   ├── users.go         # Usersコマンド
│   │   ├── newsletters.go   # Newslettersコマンド
│   │   ├── tiers.go         # Tiersコマンド
│   │   ├── offers.go        # Offersコマンド
│   │   ├── themes.go        # Themesコマンド
│   │   └── webhooks.go      # Webhooksコマンド
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
│   │   ├── images_test.go
│   │   ├── members.go       # Members API
│   │   ├── members_test.go
│   │   ├── users.go         # Users API
│   │   ├── users_test.go
│   │   ├── newsletters.go   # Newsletters API
│   │   ├── newsletters_test.go
│   │   ├── tiers.go         # Tiers API
│   │   ├── tiers_test.go
│   │   ├── offers.go        # Offers API
│   │   ├── offers_test.go
│   │   ├── themes.go        # Themes API
│   │   ├── themes_test.go
│   │   ├── webhooks.go      # Webhooks API
│   │   └── webhooks_test.go
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
- `internal/ghostapi/` - APIクライアント（47テスト）
  - `client.go`, `jwt.go` - 11テスト
  - `posts.go` - 7テスト
  - `pages.go` - 5テスト
  - `tags.go` - 7テスト
  - `images.go` - 2テスト
  - `members.go` - 6テスト
  - `users.go` - 7テスト
  - `newsletters.go` - 6テスト（List、Get、Create、Update）
  - `tiers.go` - 6テスト（List、Get、Create、Update）
  - `offers.go` - 6テスト（List、Get、Create、Update）
  - `themes.go` - 3テスト
  - `webhooks.go` - 3テスト
- `internal/outfmt/` - 出力フォーマット（5テスト）

合計: 66テスト、すべてパス

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

Phase 7が完了し、主要なGhost Admin API機能の実装がすべて完了しました。

今後の拡張機能については `docs/NEXT_STEPS.md` を参照してください。
