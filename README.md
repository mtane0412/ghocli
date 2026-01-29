# gho - Ghost Admin API CLI

gog-cliの使用感を備えたGhost Admin APIのCLIツールです。

## 特徴

- Ghost Admin APIの操作をコマンドラインから実行
- OSのキーリング（macOS Keychain、Linux Secret Service、Windows Credential Manager）によるAPIキーの安全な保存
- マルチサイト対応（エイリアス機能）
- JSON/テーブル/TSV形式での出力

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
