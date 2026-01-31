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
gho posts info <id-or-slug>

# 投稿の本文コンテンツを表示
gho posts cat <id-or-slug>
gho posts cat <id-or-slug> --format text    # テキスト形式で表示
gho posts cat <id-or-slug> --format lexical # Lexical JSON形式で表示

# 新規投稿を作成
gho posts create --title "タイトル" --html "本文" --status draft

# 投稿を更新
gho posts update <id> --title "新しいタイトル"
gho posts update <id> --html "新しい本文"

# 投稿を削除
gho posts delete <id>

# 投稿を公開
gho posts publish <id>

# 投稿をコピー（新しい下書きとして作成）
gho posts copy <id-or-slug>
gho posts copy <id-or-slug> --title "新しいタイトル"
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
gho pages info <id-or-slug>

# ページの本文コンテンツを表示
gho pages cat <id-or-slug>
gho pages cat <id-or-slug> --format text    # テキスト形式で表示
gho pages cat <id-or-slug> --format lexical # Lexical JSON形式で表示

# 新規ページを作成
gho pages create --title "タイトル" --html "本文"

# ページを更新
gho pages update <id> --title "新しいタイトル"
gho pages update <id> --html "新しい本文"

# ページを削除
gho pages delete <id>

# ページをコピー（新しい下書きとして作成）
gho pages copy <id-or-slug>
gho pages copy <id-or-slug> --title "新しいタイトル"
```

### Tags（タグ）

```bash
# タグ一覧を取得
gho tags list

# 件数を制限
gho tags list --limit 10

# タグ詳細を取得（IDまたはSlugで指定）
gho tags info <id-or-slug>
gho tags info slug:technology

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
gho members info <id>

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
gho users info <id-or-slug>
gho users info slug:john-doe

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
gho newsletters info <id-or-slug>
gho newsletters info slug:weekly-newsletter

# 新規ニュースレターを作成
gho newsletters create --name "週刊ニュースレター" --description "毎週金曜日配信"

# 送信者情報を指定して作成
gho newsletters create --name "月刊レター" --sender-name "編集部" --sender-email "editor@example.com"

# ニュースレターを更新
gho newsletters update <id> --name "新しい名前"
gho newsletters update <id> --visibility paid --subscribe-on-signup=false
```

### Tiers（ティア）

```bash
# ティア一覧を取得
gho tiers list

# 価格情報を含めて取得
gho tiers list --include monthly_price,yearly_price

# ティア詳細を取得（IDまたはSlugで指定）
gho tiers info <id-or-slug>
gho tiers info slug:premium

# 新規ティアを作成（無料プラン）
gho tiers create --name "フリープラン" --type free

# 有料ティアを作成
gho tiers create --name "プレミアム" --type paid --monthly-price 1000 --yearly-price 10000 --currency JPY

# 特典付きティアを作成
gho tiers create --name "VIP" --type paid --monthly-price 3000 --benefits "優先サポート" --benefits "限定コンテンツ"

# ティアを更新
gho tiers update <id> --name "新プレミアム"
gho tiers update <id> --monthly-price 1200 --yearly-price 12000
```

### Offers（オファー）

```bash
# オファー一覧を取得
gho offers list

# フィルターを適用
gho offers list --filter "status:active"

# オファー詳細を取得
gho offers info <id>

# パーセント割引のオファーを作成
gho offers create --name "新規会員割引" --code "WELCOME2024" --type percent --amount 20 --tier-id <tier-id>

# 固定金額割引のオファーを作成
gho offers create --name "500円オフ" --code "SAVE500" --type fixed --amount 500 --currency JPY --tier-id <tier-id>

# 期間限定オファーを作成
gho offers create --name "3ヶ月割引" --code "TRIAL3M" --type percent --amount 50 --duration repeating --duration-in-months 3 --tier-id <tier-id>

# オファーを更新
gho offers update <id> --name "新規登録キャンペーン"
gho offers update <id> --amount 30
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

## 出力形式

ghoは3つの出力形式をサポートしています：

### テーブル形式（デフォルト）

人間が読みやすい形式で出力します。

**info系コマンド（単一アイテム）**:
```bash
$ gho site
title        はなしのタネ
description  技術・学問・ゲーム・田舎暮らしを中心に...
url          https://hanatane.net/
version      6.8
```

**list系コマンド（複数アイテム）**:
```bash
$ gho posts list --limit 3
ID                        TITLE                               STATUS     CREATED     PUBLISHED
697b61d44921c40001f01aa3  CLIを使えない/使わない              draft      2026-01-29
696ce7244921c40001f017ed  非エンジニアおじさんの開発環境2026  published  2026-01-18  2026-01-28
```

### Plain形式（TSV）

スクリプトやパイプラインでの処理に適したタブ区切り形式です。

```bash
$ gho site --plain
title	はなしのタネ
description	技術・学問・ゲーム・田舎暮らしを中心に...
url	https://hanatane.net/
version	6.8

$ gho posts list --plain --limit 2
ID	TITLE	STATUS	CREATED	PUBLISHED
697b61d44921c40001f01aa3	CLIを使えない/使わない	draft	2026-01-29
696ce7244921c40001f017ed	非エンジニアおじさんの開発環境2026	published	2026-01-18	2026-01-28
```

### JSON形式

プログラムからの処理やAPI連携に適した形式です。

```bash
$ gho site --json
{
  "site": {
    "title": "はなしのタネ",
    "description": "技術・学問・ゲーム・田舎暮らしを中心に...",
    "url": "https://hanatane.net/",
    "version": "6.8"
  }
}
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

# 確認をスキップ（作成/更新/削除コマンド）
gho posts delete <id> --force
gho newsletters create --name "Test" --force
gho tiers update <id> --name "New Name" --force

# 詳細ログを表示
gho -v posts list
```

## アーキテクチャ

ghoはgogcliの設計パターンに基づいて実装されており、以下の特徴を持ちます：

### 設計原則

- **Context伝搬**: すべてのコマンドでcontextを使用し、出力モードやUIインスタンスを安全に伝搬
- **終了コード管理**: ExitError型による適切な終了コード制御
- **TDD**: テスト駆動開発による堅牢な実装
- **型安全性**: Goの型システムを最大限活用

### 主要コンポーネント

- **internal/cmd**: コマンド実装（全コマンドが`Run(ctx context.Context, root *RootFlags) error`シグネチャ）
- **internal/outfmt**: 出力フォーマット管理（JSON/Table/Plain）
- **internal/ui**: UI出力管理（stdout/stderr分離）
- **internal/ghostapi**: Ghost Admin API クライアント
- **internal/secrets**: OSキーリング統合
- **internal/config**: 設定ファイル管理

詳細は以下のドキュメントを参照してください：
- [設計統一の進捗状況](./docs/gogcli-alignment-status.md)
- [残りタスクの実装ガイド](./docs/remaining-tasks-guide.md)

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
