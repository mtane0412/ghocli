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
   - テーブル形式（人間向け、gogcliスタイル）
   - Plain形式（TSV、プログラム連携向け）

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
   gho posts info <id-or-slug>   # 旧: gho posts get（後方互換あり）
   gho posts create --title "..." [--html "..."] [--status draft|published]
   gho posts update <id> [--title "..."] [--html "..."]
   gho posts delete <id>
   gho posts publish <id>
   ```

4. **Pagesコマンド** (`internal/cmd/pages.go`)
   ```
   gho pages list [--status draft|published|scheduled] [--limit N]
   gho pages info <id-or-slug>   # 旧: gho pages get（後方互換あり）
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
   gho tags info <id-or-slug>       # 旧: gho tags get（後方互換あり）
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
   gho members info <id>            # 旧: gho members get（後方互換あり）
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
   gho users info <id-or-slug>      # 旧: gho users get（後方互換あり）
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
   gho newsletters info <id-or-slug>   # 旧: gho newsletters get（後方互換あり）
   gho newsletters create --name "..." [--description "..."] [--visibility members|paid]
   gho newsletters update <id> [--name "..."] [--visibility "..."] [--sender-name "..."]
   ```

5. **Tiersコマンド** (`internal/cmd/tiers.go`)
   ```
   gho tiers list [--limit N] [--page N] [--include monthly_price,yearly_price]
   gho tiers info <id-or-slug>         # 旧: gho tiers get（後方互換あり）
   gho tiers create --name "..." [--type free|paid] [--monthly-price N] [--yearly-price N]
   gho tiers update <id> [--name "..."] [--monthly-price N] [--yearly-price N]
   ```

6. **Offersコマンド** (`internal/cmd/offers.go`)
   ```
   gho offers list [--limit N] [--page N] [--filter "..."]
   gho offers info <id>                # 旧: gho offers get（後方互換あり）
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
github.com/k3a/html2text v1.3.0           # HTML→テキスト変換
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

## バグ修正履歴

### 2026-01-30: JWT署名エラーとPosts/Pages更新時の編集ロックエラーの修正

**問題**:
1. Ghost Admin APIとの通信時に`Invalid token: invalid signature`エラーが発生
2. 投稿・ページの更新時に`Someone else is editing this post`エラーが発生

**原因**:
1. Ghost Admin APIのシークレットキーは16進数文字列として提供されるが、JWT署名時にバイナリにデコードせずに直接使用していた
2. 投稿・ページの更新時に、サーバーから取得した元の`updated_at`タイムスタンプではなく、新しいタイムスタンプを生成して送信していた（Ghost APIの楽観的ロック機構に違反）

**修正箇所**:
- `internal/ghostapi/jwt.go:46-50` - シークレットを16進数からバイナリにデコードしてから署名
- `internal/ghostapi/jwt_test.go:9,58-63` - テストコードも16進数デコードに対応
- 全テストファイル - テスト用のシークレットを16進数形式に統一
- `internal/cmd/posts.go:202` - `UpdatedAt: time.Now()`を`UpdatedAt: existingPost.UpdatedAt`に修正
- `internal/cmd/posts.go:310` - publishコマンドでも同様の修正
- `internal/cmd/pages.go:201` - ページ更新でも同様の修正

**テスト追加**:
- `TestUpdatePost_updated_atを保持して更新` - 更新時に元の`updated_at`を送信することを確認するテスト

**動作確認**:
- すべての読み取り操作（Posts、Tags、Users、Newsletters、Tiers、Pagesなど）が正常に動作
- 書き込み操作（作成・更新・削除）が正常に動作
- 81テストすべてがパス

**参考**:
- [Ghost Admin API Overview](https://docs.ghost.org/admin-api)
- [Bash Example of Ghost JWT Auth](https://gist.github.com/ErisDS/6334f0e70ec7390ec08530d5ef9bd0d5)

### ✅ Phase 8: コマンド設計改善（進行中）

**開始日**: 2026-01-31

**目的**: gogcliのコマンド設計パターンを参考に、ghoのコマンド体系を改善

**実装内容**:

#### 8.1: get → info リネーム（完了）
- すべてのリソース（posts, pages, members, tags, users, newsletters, tiers, offers）で`get`コマンドを`info`にリネーム
- 後方互換性のため`get`はエイリアス（非推奨）として維持
- 非推奨メッセージをヘルプに追加

**変更例**:
```bash
# 新しいコマンド
gho posts info <id>
gho pages info <slug>
gho members info <id>

# 旧コマンド（非推奨警告付きで動作）
gho posts get <id>
```

#### 8.2: catコマンドの追加（完了）
- posts/pagesに本文コンテンツを標準出力に表示する`cat`コマンドを追加
- `--format`オプションでhtml/text/lexical形式を選択可能
- textフォーマットはk3a/html2textライブラリを使用してHTMLからプレーンテキストに変換

**使用例**:
```bash
gho posts cat <id>                      # HTML形式で出力
gho posts cat <id> --format text        # テキスト形式で出力（HTMLタグを除去）
gho posts cat <id> --format lexical     # Lexical JSON形式で出力
gho pages cat <slug> --format html      # ページの本文をHTML形式で出力
```

#### 8.3: copyコマンドの追加（完了）
- posts/pagesをコピーする`copy`コマンドを追加
- ID/UUID/Slug/URL/日時を除外して新規作成
- ステータスは常に`draft`で作成
- `--title`オプションで新しいタイトルを指定（省略時は「元タイトル (Copy)」）

**使用例**:
```bash
gho posts copy <id-or-slug>                      # 投稿をコピー（タイトルは「元タイトル (Copy)」）
gho posts copy <id-or-slug> --title "新タイトル" # カスタムタイトルでコピー
gho pages copy <slug>                            # ページをコピー
gho pages copy <slug> --title "新タイトル"       # カスタムタイトルでコピー
```

#### 8.5: 出力形式のgogcliスタイル完全移行（完了）
- info系コマンドのキー名をsnake_case（小文字）に変更
- list系コマンドのセパレーター行を削除
- テーブル形式とplain形式の統一（tabwriterによる自動整列）
- すべてのinfo系コマンドでPrintKeyValueを使用

**変更内容**:
- info系コマンド: ヘッダーなし、キー名小文字、タブ区切り
- list系コマンド: ヘッダー大文字、セパレーターなし、タブ区切り
- members/users/tags infoコマンドをPrintTableからPrintKeyValueに変更

**出力例**:
```bash
# info系（テーブル形式）
$ gho site
title        はなしのタネ
description  技術・学問・ゲーム・田舎暮らしを中心に...
url          https://hanatane.net/
version      6.8

# info系（plain形式）
$ gho site --plain
title	はなしのタネ
description	技術・学問・ゲーム...
url	https://hanatane.net/
version	6.8

# list系（テーブル形式）
$ gho posts list --limit 2
ID                        TITLE                               STATUS     CREATED     PUBLISHED
697b61d44921c40001f01aa3  CLIを使えない/使わない              draft      2026-01-29
696ce7244921c40001f017ed  非エンジニアおじさんの開発環境2026  published  2026-01-18  2026-01-28
```

#### 8.4: catコマンドのtextフォーマット実装（完了）
- posts/pagesの`cat`コマンドで`--format text`が正しく動作するように実装
- k3a/html2textライブラリを使用してHTMLからプレーンテキストに変換
- シェルリダイレクト（`gho posts cat <id> --format html > output.html`）でエクスポート可能なため、専用のexportコマンドは不要と判断

**実装箇所**:
- `internal/cmd/posts.go:939` - HTMLからテキストへの変換実装
- `internal/cmd/pages.go:497` - 同上
- `go.mod` - k3a/html2text v1.3.0を追加

**品質チェック（Phase 8.5完了時点）**:
- ✅ すべてのテストがパス（164テスト）
- ✅ 型チェック（`go vet`）成功
- ✅ Lint（golangci-lint）成功（0 issues）
- ✅ ビルド成功

**コミット**:
- `dec99de feat(cmd): Phase 8.1 - get → info リネーム（全リソース）`
- `a1d6f61 feat(cmd): Phase 8.2 - catコマンドの追加（posts, pages）`
- `18ab842 feat(cmd): Phase 8.3 - copyコマンドの追加（posts, pages）`
- `bf260b7 feat(cmd): Phase 8.4 - catコマンドのtextフォーマット実装`
- `6f1a8b5 feat(outfmt): info系コマンドのヘッダー削除とgogcliスタイル採用`
- `043ac98 feat(outfmt): gogcliスタイルに完全移行`

**参考**: gogcliのコマンド設計パターン（`gog docs info/cat/copy/export`）

## 次のステップ

Phase 7が完了し、主要なGhost Admin API機能の実装がすべて完了しました。
現在、Phase 8（コマンド設計改善）を進行中です。

今後の拡張機能については `docs/NEXT_STEPS.md` を参照してください。
