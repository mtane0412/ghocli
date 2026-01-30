# gho - Ghost Admin API CLI

gog-cliの使用感を備えたGhost Admin APIのCLIツールです。

## 特徴

- Ghost Admin APIの操作をコマンドラインから実行
- OSのキーリング（macOS Keychain、Linux Secret Service、Windows Credential Manager）によるAPIキーの安全な保存
- マルチサイト対応（エイリアス機能）
- JSON/テーブル/TSV形式での出力
- TDD（テスト駆動開発）による堅牢な実装

## 実装状況

✅ **Phase 1-7 完了** - 主要なGhost Admin API機能をすべて実装済み

**実装済み機能**:
- 認証管理（Auth）
- サイト情報（Site）
- 投稿管理（Posts）
- ページ管理（Pages）
- タグ管理（Tags）
- 画像管理（Images）
- メンバー管理（Members）
- ユーザー管理（Users）
- ニュースレター（Newsletters）
- ティア（Tiers）
- オファー（Offers）
- テーマ管理（Themes）
- Webhook管理（Webhooks）

詳細は [`docs/PROJECT_STATUS.md`](./docs/PROJECT_STATUS.md) を参照してください。

## インストール

```bash
# ソースからビルド
git clone https://github.com/mtane0412/gho.git
cd gho
make build

# または
go install github.com/mtane0412/gho/cmd/gho@latest
```

## 使用方法

### 認証設定

```bash
# Ghost Admin APIキーを登録
gho auth add https://your-blog.ghost.io

# 登録済みサイト一覧
gho auth list

# 認証状態確認
gho auth status

# サイト認証を削除
gho auth remove <alias>
```

### サイト情報

```bash
# サイト情報を取得
gho site

# 特定のサイトの情報を取得
gho -s myblog site

# JSON形式で出力
gho site --json
```

### Posts（投稿）

```bash
# 投稿一覧を取得
gho posts list

# ステータスでフィルタリング
gho posts list --status draft
gho posts list --status published
gho posts list --status scheduled

# 件数を制限
gho posts list --limit 10

# 投稿詳細を取得（IDまたはSlugで指定）
gho posts get <id-or-slug>

# 新規投稿を作成
gho posts create --title "タイトル" --html "本文" --status draft

# 投稿を更新
gho posts update <id> --title "新しいタイトル"
gho posts update <id> --html "新しい本文"

# 投稿を削除
gho posts delete <id>

# 投稿を公開
gho posts publish <id>
```

### Pages（固定ページ）

```bash
# ページ一覧を取得
gho pages list

# ステータスでフィルタリング
gho pages list --status draft
gho pages list --status published
gho pages list --status scheduled

# 件数を制限
gho pages list --limit 10

# ページ詳細を取得（IDまたはSlugで指定）
gho pages get <id-or-slug>

# 新規ページを作成
gho pages create --title "タイトル" --html "本文"

# ページを更新
gho pages update <id> --title "新しいタイトル"
gho pages update <id> --html "新しい本文"

# ページを削除
gho pages delete <id>
```

### Tags（タグ）

```bash
# タグ一覧を取得
gho tags list

# 件数を制限
gho tags list --limit 10

# タグ詳細を取得（IDまたはSlugで指定）
gho tags get <id-or-slug>
gho tags get slug:technology

# 新規タグを作成
gho tags create --name "Technology" --description "技術関連の記事"

# タグの可視性を指定
gho tags create --name "Internal" --visibility internal

# タグを更新
gho tags update <id> --name "Tech" --description "新しい説明"

# タグを削除
gho tags delete <id>
```

### Images（画像）

```bash
# 画像をアップロード
gho images upload path/to/image.jpg

# 用途を指定してアップロード
gho images upload avatar.png --purpose profile_image
gho images upload icon.png --purpose icon

# 参照IDを指定
gho images upload banner.jpg --ref post-123
```

### Members（メンバー）

```bash
# メンバー一覧を取得
gho members list

# 件数を制限
gho members list --limit 10

# フィルターを適用
gho members list --filter "status:paid"

# メンバー詳細を取得
gho members get <id>

# 新規メンバーを作成
gho members create --email "user@example.com" --name "山田太郎"

# ラベル付きでメンバーを作成
gho members create --email "user@example.com" --labels "VIP,Premium"

# メンバーを更新
gho members update <id> --name "田中花子" --note "重要顧客"

# メンバーを削除
gho members delete <id>
```

### Users（ユーザー）

```bash
# ユーザー一覧を取得
gho users list

# ロール情報を含めて取得
gho users list --include roles

# 投稿数を含めて取得
gho users list --include count.posts

# ユーザー詳細を取得（IDまたはSlugで指定）
gho users get <id-or-slug>
gho users get slug:john-doe

# ユーザー情報を更新
gho users update <id> --name "新しい名前" --bio "新しい自己紹介"
gho users update <id> --location "Tokyo" --website "https://example.com"
```

### Newsletters（ニュースレター）

```bash
# ニュースレター一覧を取得
gho newsletters list

# フィルターを適用
gho newsletters list --filter "status:active"

# ニュースレター詳細を取得（IDまたはSlugで指定）
gho newsletters get <id-or-slug>
gho newsletters get slug:weekly-newsletter
```

### Tiers（ティア）

```bash
# ティア一覧を取得
gho tiers list

# 価格情報を含めて取得
gho tiers list --include monthly_price,yearly_price

# ティア詳細を取得（IDまたはSlugで指定）
gho tiers get <id-or-slug>
gho tiers get slug:premium
```

### Offers（オファー）

```bash
# オファー一覧を取得
gho offers list

# フィルターを適用
gho offers list --filter "status:active"

# オファー詳細を取得
gho offers get <id>
```

### Themes（テーマ）

```bash
# テーマ一覧を取得
gho themes list

# テーマをアップロード
gho themes upload path/to/theme.zip

# テーマを有効化
gho themes activate casper
```

### Webhooks（Webhook）

```bash
# Webhookを作成
gho webhooks create --event post.published --target-url https://example.com/webhook

# 名前付きでWebhookを作成
gho webhooks create --event member.added --target-url https://example.com/webhook --name "Member notification"

# Webhookを更新
gho webhooks update <id> --target-url https://new-example.com/webhook

# Webhookを削除
gho webhooks delete <id>
```

## グローバルオプション

すべてのコマンドで以下のオプションが使用できます：

```bash
# JSON形式で出力
gho posts list --json

# TSV形式で出力（スクリプト連携向け）
gho posts list --plain

# 特定のサイトを指定
gho -s myblog posts list

# 確認をスキップ（削除コマンドなど）
gho posts delete <id> --force

# 詳細ログを表示
gho -v posts list
```

## 開発

### テスト実行

```bash
make test
```

### Lint実行

```bash
make lint
```

### 型チェック

```bash
make type-check
```

### ビルド

```bash
make build
```

## ライセンス

MIT
