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

### 🚧 Phase 2: コンテンツ管理（Posts/Pages）（未実装）

**予定内容**:

1. **Posts API** (`internal/ghostapi/posts.go`, `internal/cmd/posts.go`)
   ```
   gho posts list [--status draft|published|scheduled] [--limit N]
   gho posts get <id-or-slug>
   gho posts create --title "..." [--html "..."]
   gho posts update <id> --title "..."
   gho posts delete <id>
   gho posts publish <id>
   ```

2. **Pages API** (`internal/ghostapi/pages.go`, `internal/cmd/pages.go`)
   ```
   gho pages list
   gho pages get <id-or-slug>
   gho pages create --title "..."
   gho pages update <id> ...
   gho pages delete <id>
   ```

### 📋 Phase 3: タクソノミー + メディア（未実装）

**予定内容**:

1. **Tags API**
   ```
   gho tags list
   gho tags get <id-or-slug>
   gho tags create --name "..."
   gho tags update <id> --name "..."
   gho tags delete <id>
   ```

2. **Images API**
   ```
   gho images upload <file-path>
   ```

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
│   │   └── site.go          # サイト情報コマンド
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
│   │   └── jwt_test.go
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
- `internal/secrets/` - キーリング統合（5テスト）
- `internal/ghostapi/` - APIクライアント（9テスト）
- `internal/outfmt/` - 出力フォーマット（5テスト）

合計: 25テスト、すべてパス

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

Phase 2の実装を開始します。詳細は `docs/NEXT_STEPS.md` を参照してください。
