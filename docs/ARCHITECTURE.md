# gho アーキテクチャ設計

## 概要

ghoはGhost Admin APIのCLIツールです。gog-cliの設計パターンを参考にし、シンプルで保守性の高いアーキテクチャを採用しています。

## プロジェクト構造

```
gho/
├── cmd/gho/
│   └── main.go              # エントリーポイント
├── internal/
│   ├── cmd/                  # CLIコマンド定義
│   │   ├── root.go          # CLI構造体、RootFlags
│   │   ├── auth.go          # 認証コマンド
│   │   ├── site.go          # サイト情報コマンド
│   │   ├── posts.go         # Postsコマンド（Phase 2）
│   │   └── pages.go         # Pagesコマンド（Phase 2）
│   ├── config/              # 設定ファイル管理
│   │   ├── config.go
│   │   └── config_test.go
│   ├── secrets/             # キーリング統合
│   │   ├── store.go
│   │   └── store_test.go
│   ├── ghostapi/            # Ghost APIクライアント
│   │   ├── client.go        # HTTPクライアント + JWT生成
│   │   ├── client_test.go
│   │   ├── jwt.go           # JWT生成
│   │   ├── jwt_test.go
│   │   ├── posts.go         # Posts API（Phase 2）
│   │   └── pages.go         # Pages API（Phase 2）
│   ├── outfmt/              # 出力フォーマット
│   │   ├── outfmt.go
│   │   └── outfmt_test.go
│   └── errfmt/              # エラーフォーマット（今後実装）
├── docs/                    # ドキュメント
├── go.mod
├── go.sum
├── Makefile
├── .golangci.yml
├── .gitignore
└── README.md
```

## レイヤー構成

```
┌─────────────────────────────────────┐
│          CLI Layer (cmd/)           │  ← ユーザーインターフェース
│  - コマンド定義                      │
│  - フラグパース                      │
│  - 入力検証                          │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│      Business Logic Layer           │
│  - config/  : 設定管理               │  ← ビジネスロジック
│  - secrets/ : 認証情報管理           │
│  - ghostapi/: API操作                │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│      Infrastructure Layer           │
│  - outfmt/  : 出力フォーマット       │  ← インフラストラクチャ
│  - errfmt/  : エラーフォーマット     │
│  - HTTP Client                       │
│  - OS Keyring                        │
└─────────────────────────────────────┘
```

## コンポーネント設計

### 1. CLI Layer (`internal/cmd/`)

**責務**: ユーザーからの入力を受け取り、適切なビジネスロジックを呼び出す

**主要コンポーネント**:

- **RootFlags**: すべてのコマンドで共通のフラグ
  ```go
  type RootFlags struct {
      Site    string // サイトエイリアスまたはURL
      JSON    bool   // JSON形式で出力
      Plain   bool   // TSV形式で出力
      Force   bool   // 確認をスキップ
      Verbose bool   // 詳細ログを有効化
  }
  ```

- **CLI**: Kongで定義されるCLI構造体
  ```go
  type CLI struct {
      RootFlags `embed:""`
      Version   kong.VersionFlag
      Auth      AuthCmd
      Site      SiteCmd
      Posts     PostsCmd
      Pages     PagesCmd
      // ...
  }
  ```

**設計パターン**: Command Pattern（Kongが内部的に使用）

### 2. Config Layer (`internal/config/`)

**責務**: 設定ファイルの読み書き、サイト管理

**主要機能**:

- 設定ファイルパス: `~/.config/gho/config.json`
- マルチサイト対応（エイリアス機能）
- デフォルトサイト管理

**設定ファイル形式**:
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

**主要メソッド**:
- `Load(path string) (*Config, error)` - 設定をロード
- `Save(path string) error` - 設定を保存
- `AddSite(alias, url string)` - サイトを追加
- `GetSiteURL(aliasOrURL string) (string, bool)` - URLを取得

### 3. Secrets Layer (`internal/secrets/`)

**責務**: Admin APIキーの安全な保存・取得

**キーリングバックエンド**:
- macOS: Keychain
- Linux: Secret Service (GNOME Keyring, KWallet)
- Windows: Credential Manager
- Fallback: 暗号化ファイル

**主要メソッド**:
- `NewStore(backend, fileDir string) (*Store, error)` - ストア作成
- `Set(alias, apiKey string) error` - APIキー保存
- `Get(alias string) (string, error)` - APIキー取得
- `Delete(alias string) error` - APIキー削除
- `List() ([]string, error)` - 保存済みエイリアス一覧
- `ParseAdminAPIKey(apiKey string) (id, secret string, err error)` - APIキーパース

**セキュリティ**:
- APIキーはOSキーリングに保存（プレーンテキストでファイルに保存しない）
- Fallback（fileバックエンド）の場合もパスワード保護

### 4. Ghost API Layer (`internal/ghostapi/`)

**責務**: Ghost Admin APIとの通信

**主要コンポーネント**:

#### Client
HTTPクライアントとJWT生成を統合

```go
type Client struct {
    baseURL    string
    keyID      string
    secret     string
    httpClient *http.Client
}
```

**主要メソッド**:
- `NewClient(baseURL, keyID, secret string) (*Client, error)`
- `doRequest(method, path string, body io.Reader) ([]byte, error)`
- `GetSite() (*Site, error)`

#### JWT生成
Ghost Admin APIはHS256で署名されたJWTを要求

```go
func GenerateJWT(keyID, secret string) (string, error)
```

**JWTクレーム**:
```json
{
  "iat": 1234567890,      // 発行時刻（Unix時間）
  "exp": 1234568190,      // 有効期限（iat + 5分）
  "aud": "/admin/"        // Ghost Admin APIのパス
}
```

**JWTヘッダー**:
```json
{
  "alg": "HS256",         // 署名アルゴリズム
  "typ": "JWT",
  "kid": "64fac5417..."   // APIキーID
}
```

#### API型定義

各APIリソースに対応する型を定義

```go
type Site struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    URL         string `json:"url"`
    Version     string `json:"version"`
}

type Post struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Slug        string     `json:"slug"`
    HTML        string     `json:"html,omitempty"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    PublishedAt *time.Time `json:"published_at,omitempty"`
    Tags        []Tag      `json:"tags,omitempty"`
    Authors     []Author   `json:"authors,omitempty"`
}
```

### 5. Output Format Layer (`internal/outfmt/`)

**責務**: 出力フォーマットの統一管理

**サポート形式**:

| モード | フラグ | 用途 | 形式 |
|--------|--------|------|------|
| Table | (default) | 人間向け | カラム揃え、ヘッダー付き |
| JSON | `--json` | プログラム連携 | JSON形式 |
| Plain | `--plain` | パイプ処理 | TSV形式 |

**主要メソッド**:
- `NewFormatter(writer io.Writer, mode string) *Formatter`
- `Print(data interface{}) error` - 任意のデータを出力
- `PrintTable(headers []string, rows [][]string) error` - テーブル出力
- `PrintMessage(message string)` - メッセージ出力
- `PrintError(message string)` - エラー出力

**出力例**:

**Table形式**:
```
Alias   URL                           Default
------  ----------------------------  -------
myblog  https://myblog.ghost.io       *
work    https://blog.company.com
```

**JSON形式**:
```json
[
  {
    "Alias": "myblog",
    "URL": "https://myblog.ghost.io",
    "Default": "*"
  },
  {
    "Alias": "work",
    "URL": "https://blog.company.com",
    "Default": ""
  }
]
```

**Plain形式（TSV）**:
```
Alias	URL	Default
myblog	https://myblog.ghost.io	*
work	https://blog.company.com
```

## 認証フロー

```
1. ユーザーがGhost Adminで Custom Integration を作成
   ↓
2. `gho auth add https://myblog.ghost.io` を実行
   ↓
3. APIキー（id:secret形式）を入力
   ↓
4. APIキーをパース（secrets.ParseAdminAPIKey）
   ↓
5. `/ghost/api/admin/site/` で検証
   - JWTを生成（jwt.GenerateJWT）
   - HTTPリクエストを実行（client.GetSite）
   ↓
6. キーリングに保存（secrets.Store.Set）
   ↓
7. 設定ファイルにサイトを追加（config.Config.AddSite）
   ↓
8. 設定ファイルを保存（config.Config.Save）
```

## APIリクエストフロー

```
1. ユーザーがコマンドを実行（例: gho site）
   ↓
2. RootFlagsからサイトを決定
   - -s フラグで指定されたサイト
   - または設定のdefault_site
   ↓
3. 設定ファイルからURLを取得（config.Config.GetSiteURL）
   ↓
4. キーリングからAPIキーを取得（secrets.Store.Get）
   ↓
5. APIキーをパース（secrets.ParseAdminAPIKey）
   ↓
6. APIクライアントを作成（ghostapi.NewClient）
   ↓
7. JWTを生成（ghostapi.GenerateJWT）
   ↓
8. HTTPリクエストを実行（ghostapi.Client.doRequest）
   - Authorization: Ghost <JWT>
   - Accept: application/json
   ↓
9. レスポンスをパース
   ↓
10. 出力フォーマットで表示（outfmt.Formatter）
```

## エラーハンドリング

### エラーの種類

1. **設定エラー**
   - 設定ファイルが見つからない
   - サイトが登録されていない
   - デフォルトサイトが設定されていない

2. **認証エラー**
   - APIキーが無効
   - APIキーの形式が不正
   - キーリングへのアクセスエラー

3. **APIエラー**
   - HTTP エラー（401, 404, 500など）
   - レスポンスのパースエラー
   - ネットワークエラー

4. **入力エラー**
   - 必須パラメータが不足
   - パラメータの形式が不正

### エラーメッセージの設計

すべてのエラーメッセージは以下の形式：

```
Error: <エラーの説明>
```

例:
```
Error: サイト 'myblog' が見つかりません
Error: APIキーの検証に失敗: Unauthorized
Error: 設定の読み込みに失敗: open /Users/user/.config/gho/config.json: no such file or directory
```

## テスト戦略

### ユニットテスト

各コンポーネントは独立してテスト可能：

- **config**: 設定ファイルの読み書き
- **secrets**: キーリング操作（fileバックエンドでテスト）
- **ghostapi**: HTTPクライアント（httptestでモック）
- **outfmt**: 出力フォーマット（bytes.Bufferで検証）

### テストカバレッジ目標

- コアロジック: 80%以上
- API層: 70%以上
- CLI層: 60%以上（手動テストでカバー）

## パフォーマンス考慮事項

### JWT生成

- 各APIリクエストごとにJWTを生成
- 有効期限は5分（Ghost Admin APIの要件）
- キャッシュは不要（生成コストは低い）

### HTTP接続

- タイムアウト: 30秒
- Keep-Alive: デフォルト有効
- 複数リクエストでは接続を再利用

### キーリングアクセス

- 初回アクセス時にキーリングをオープン
- 複数の操作で再利用可能
- パスワード入力はバックエンドに依存

## セキュリティ考慮事項

### APIキーの保存

- OSキーリングに保存（プレーンテキストでファイルに保存しない）
- 設定ファイルにはURLのみ保存（APIキーは含まない）
- ファイルバックエンドはパスワード保護

### JWT

- 有効期限5分（短命）
- HS256署名（Ghost Admin APIの要件）
- ヘッダーにkid（キーID）を含む

### ファイルパーミッション

- 設定ファイル: 0600（所有者のみ読み書き）
- キーリングファイル: 0600

## 拡張性

### 新しいAPIリソースの追加

1. `internal/ghostapi/` に型定義を追加
2. `internal/ghostapi/` にAPI関数を追加
3. `internal/cmd/` にコマンドを追加
4. テストを追加
5. ドキュメントを更新

### 新しい出力形式の追加

1. `internal/outfmt/` に新しいフォーマッターを追加
2. `RootFlags` に新しいフラグを追加
3. `GetOutputMode()` でモードを返すように修正

### 新しいキーリングバックエンドの追加

99designs/keyringがサポートするバックエンドは自動的に利用可能
