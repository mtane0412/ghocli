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
