# 次のステップ

## 現在の状態

✅ **Phase 1: 基盤構築** - 完了（2026-01-29）

- 設定システム
- キーリング統合
- Ghost APIクライアント（JWT生成、HTTPクライアント）
- 出力フォーマット
- 認証コマンド（auth add/list/remove/status）
- サイト情報取得コマンド（site）

## Phase 2: コンテンツ管理（Posts/Pages）

### 目標

Posts/Pagesの作成、更新、削除、公開機能を実装する

### タスクリスト

#### 1. Posts API実装

- [ ] `internal/ghostapi/posts.go` を作成
  - [ ] Post型定義
  - [ ] ListOptions型定義（limit, status, filterなど）
  - [ ] テスト作成（`posts_test.go`）
  - [ ] `ListPosts(options ListOptions) ([]Post, error)` 実装
  - [ ] `GetPost(idOrSlug string) (*Post, error)` 実装
  - [ ] `CreatePost(post *Post) (*Post, error)` 実装
  - [ ] `UpdatePost(id string, post *Post) (*Post, error)` 実装
  - [ ] `DeletePost(id string) error` 実装

#### 2. Postsコマンド実装

- [ ] `internal/cmd/posts.go` を作成
  - [ ] PostsCmd構造体定義
  - [ ] PostsListCmd実装
    - [ ] `--status` フラグ（draft/published/scheduled）
    - [ ] `--limit` フラグ
  - [ ] PostsGetCmd実装
  - [ ] PostsCreateCmd実装
    - [ ] `--title` フラグ（必須）
    - [ ] `--html` フラグ
    - [ ] `--status` フラグ
  - [ ] PostsUpdateCmd実装
  - [ ] PostsDeleteCmd実装
  - [ ] PostsPublishCmd実装

#### 3. Pages API実装

- [ ] `internal/ghostapi/pages.go` を作成
  - [ ] Page型定義
  - [ ] テスト作成（`pages_test.go`）
  - [ ] `ListPages(options ListOptions) ([]Page, error)` 実装
  - [ ] `GetPage(idOrSlug string) (*Page, error)` 実装
  - [ ] `CreatePage(page *Page) (*Page, error)` 実装
  - [ ] `UpdatePage(id string, page *Page) (*Page, error)` 実装
  - [ ] `DeletePage(id string) error` 実装

#### 4. Pagesコマンド実装

- [ ] `internal/cmd/pages.go` を作成
  - [ ] PagesCmd構造体定義
  - [ ] PagesListCmd実装
  - [ ] PagesGetCmd実装
  - [ ] PagesCreateCmd実装
  - [ ] PagesUpdateCmd実装
  - [ ] PagesDeleteCmd実装

#### 5. CLIに統合

- [ ] `internal/cmd/root.go` に PostsCmd と PagesCmd を追加

#### 6. 品質チェック

- [ ] すべてのテストがパス（`make test`）
- [ ] 型チェック成功（`make type-check`）
- [ ] Lint成功（`make lint`）
- [ ] ビルド成功（`make build`）

#### 7. 動作確認

- [ ] `./gho posts list` でPosts一覧が表示される
- [ ] `./gho posts get <slug>` で投稿詳細が表示される
- [ ] `./gho posts create --title "Test" --status draft` で投稿が作成される
- [ ] `./gho posts update <id> --title "Updated"` で投稿が更新される
- [ ] `./gho posts publish <id>` で投稿が公開される
- [ ] `./gho posts delete <id>` で投稿が削除される
- [ ] `./gho pages list` でPages一覧が表示される
- [ ] `./gho pages create --title "Test Page"` でページが作成される

#### 8. ドキュメント更新

- [ ] `docs/PROJECT_STATUS.md` を更新
- [ ] `README.md` にPosts/Pagesコマンドの使用例を追加

#### 9. コミット

- [ ] Phase 2完了のコミットを作成

### 実装の開始方法

```bash
# featureブランチを作成
git checkout -b feature/phase2-content-management

# Posts APIのテストから開始（TDD）
# internal/ghostapi/posts_test.go を作成
```

### 参考: Ghost Admin API仕様

**Posts API**:
- エンドポイント: `/ghost/api/admin/posts/`
- メソッド: GET, POST, PUT, DELETE
- パラメータ: `limit`, `status`, `filter`, `include`

**Pages API**:
- エンドポイント: `/ghost/api/admin/pages/`
- メソッド: GET, POST, PUT, DELETE

詳細: https://ghost.org/docs/admin-api/

### 実装時の注意点

1. **TDD原則を厳守**
   - テストを先に書く（RED）
   - 最小限の実装（GREEN）
   - リファクタリング（REFACTOR）

2. **エラーハンドリング**
   - エラーメッセージは日本語で具体的に
   - エラーのラップに `fmt.Errorf` と `%w` を使用

3. **出力フォーマット**
   - JSON/Table/Plain形式をサポート
   - RootFlags.GetOutputMode() で形式を決定

4. **コードコメント**
   - ファイル冒頭にコメント
   - 関数に詳細なコメント
   - 複雑な処理には日本語コメント

5. **Git ワークフロー**
   - mainブランチへの直接コミット禁止
   - featureブランチで作業
   - コミット前に品質チェック

## Phase 3以降の予定

### Phase 3: タクソノミー + メディア

- Tags API（list/get/create/update/delete）
- Images API（upload）

### Phase 4: Members管理

- Members API（list/get/create/update/delete）

### Phase 5: Users管理

- Users API（list/get/update）

### Phase 6: Newsletters/Tiers/Offers

- Newsletters API（list/get）
- Tiers API（list/get）
- Offers API（list/get）

### Phase 7: Themes/Webhooks

- Themes API（list/upload/activate/delete）
- Webhooks API（list/create/delete）

## 質問・相談

実装中に疑問が生じた場合：

1. `docs/ARCHITECTURE.md` でアーキテクチャを確認
2. `docs/DEVELOPMENT_GUIDE.md` で開発ガイドを確認
3. Phase 1の実装を参考にする
4. Ghost Admin APIドキュメントを参照

## フィードバック

実装完了後：

1. `docs/PROJECT_STATUS.md` を更新
2. `docs/NEXT_STEPS.md` を更新（Phase 3に移行）
3. 学んだことや改善点を記録
