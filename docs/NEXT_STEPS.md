# 次のステップ

## 現在の状態

✅ **Phase 1: 基盤構築** - 完了（2026-01-29）

- 設定システム
- キーリング統合
- Ghost APIクライアント（JWT生成、HTTPクライアント）
- 出力フォーマット
- 認証コマンド（auth add/list/remove/status）
- サイト情報取得コマンド（site）

✅ **Phase 2: コンテンツ管理（Posts/Pages）** - 完了（2026-01-29）

- Posts API（ListPosts、GetPost、CreatePost、UpdatePost、DeletePost）
- Pages API（ListPages、GetPage、CreatePage、UpdatePage、DeletePage）
- Postsコマンド（list、get、create、update、delete、publish）
- Pagesコマンド（list、get、create、update、delete）

✅ **Phase 3: タクソノミー + メディア** - 完了（2026-01-30）

- Tags API（ListTags、GetTag、CreateTag、UpdateTag、DeleteTag）
- Images API（UploadImage）
- Tagsコマンド（list、get、create、update、delete）
- Imagesコマンド（upload）

✅ **Phase 4: Members管理** - 完了（2026-01-30）

- Members API（ListMembers、GetMember、CreateMember、UpdateMember、DeleteMember）
- Membersコマンド（list、get、create、update、delete）

## Phase 5: Users管理

### 目標

Users（管理者・編集者）の管理機能を実装し、Ghost Admin APIのユーザー管理機能を完成させる

### タスクリスト

#### 1. Users API実装

- [ ] `internal/ghostapi/users.go` を作成
  - [ ] User型定義（ID、Name、Slug、Email、Rolesなど）
  - [ ] UserListOptions型定義（Limit、Filter、Includeなど）
  - [ ] テスト作成（`users_test.go`）
  - [ ] `ListUsers(options UserListOptions) (*UserListResponse, error)` 実装
  - [ ] `GetUser(idOrSlug string) (*User, error)` 実装
  - [ ] `UpdateUser(id string, user *User) (*User, error)` 実装（注: Ghost APIはユーザー作成・削除をサポートしていない）

#### 2. Usersコマンド実装

- [ ] `internal/cmd/users.go` を作成
  - [ ] UsersCmd構造体定義
  - [ ] UsersListCmd実装
    - [ ] `--limit` フラグ
    - [ ] `--include` フラグ（roles、count.postsなど）
  - [ ] UsersGetCmd実装
  - [ ] UsersUpdateCmd実装
    - [ ] `--name` フラグ
    - [ ] `--slug` フラグ
    - [ ] `--bio` フラグ

#### 3. CLIに統合

- [ ] `internal/cmd/root.go` に UsersCmd を追加

#### 4. 品質チェック

- [ ] すべてのテストがパス（`make test`）
- [ ] 型チェック成功（`make type-check`）
- [ ] Lint成功（`make lint`）
- [ ] ビルド成功（`make build`）

#### 5. 動作確認

- [ ] `./gho users list` でユーザー一覧が表示される
- [ ] `./gho users get <id>` でユーザー詳細が表示される
- [ ] `./gho users update <id> --name "新しい名前"` でユーザーが更新される

#### 6. ドキュメント更新

- [ ] `docs/PROJECT_STATUS.md` を更新
- [ ] `docs/NEXT_STEPS.md` を更新（Phase 6に移行）

#### 7. コミット

- [ ] Phase 5完了のコミットを作成

### 実装の開始方法

```bash
# featureブランチを作成
git checkout -b feature/phase5-users

# Users APIのテストから開始（TDD）
# internal/ghostapi/users_test.go を作成
```

### 参考: Ghost Admin API仕様

**Users API**:
- エンドポイント: `/ghost/api/admin/users/`
- メソッド: GET, PUT（注: POST/DELETEは利用不可）
- パラメータ: `limit`, `filter`, `include` (roles, count.postsなど)
- 主要フィールド: `id`, `name`, `slug`, `email`, `roles`, `bio`, `location`
- スラッグ指定: `slug:user-slug` 形式でスラッグ検索が可能

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

6. **Users特有の注意点**
   - Ghost APIはユーザーの作成・削除をサポートしていない（ダッシュボードからのみ）
   - rolesフィールドは読み取り専用
   - スラッグによる検索に対応（Tags APIと同じパターン）

## Phase 6以降の予定

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
2. `docs/NEXT_STEPS.md` を更新（Phase 6に移行）
3. 学んだことや改善点を記録
