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

## Phase 4: Members管理

### 目標

Members（購読者）の管理機能を実装し、Ghost Admin APIのメンバー管理機能を完成させる

### タスクリスト

#### 1. Members API実装

- [ ] `internal/ghostapi/members.go` を作成
  - [ ] Member型定義（ID、Email、Name、Status、CreatedAt、UpdatedAtなどGhost Admin APIのMember構造に準拠）
  - [ ] MemberListOptions型定義（Limit、Filter、Orderなど）
  - [ ] テスト作成（`members_test.go`）
  - [ ] `ListMembers(options MemberListOptions) (*MemberListResponse, error)` 実装
  - [ ] `GetMember(id string) (*Member, error)` 実装
  - [ ] `CreateMember(member *Member) (*Member, error)` 実装
  - [ ] `UpdateMember(id string, member *Member) (*Member, error)` 実装
  - [ ] `DeleteMember(id string) error` 実装

#### 2. Membersコマンド実装

- [ ] `internal/cmd/members.go` を作成
  - [ ] MembersCmd構造体定義
  - [ ] MembersListCmd実装
    - [ ] `--limit` フラグ
    - [ ] `--filter` フラグ（オプション）
    - [ ] `--order` フラグ（オプション）
  - [ ] MembersGetCmd実装
  - [ ] MembersCreateCmd実装
    - [ ] `--email` フラグ（必須）
    - [ ] `--name` フラグ（オプション）
    - [ ] `--note` フラグ（オプション）
  - [ ] MembersUpdateCmd実装
  - [ ] MembersDeleteCmd実装

#### 3. CLIに統合

- [ ] `internal/cmd/root.go` に MembersCmd を追加

#### 4. 品質チェック

- [ ] すべてのテストがパス（`make test`）
- [ ] 型チェック成功（`make type-check`）
- [ ] Lint成功（`make lint`）
- [ ] ビルド成功（`make build`）

#### 5. 動作確認

- [ ] `./gho members list` でMembers一覧が表示される
- [ ] `./gho members get <id>` でメンバー詳細が表示される
- [ ] `./gho members create --email "user@example.com" --name "山田太郎"` でメンバーが作成される
- [ ] `./gho members update <id> --name "田中太郎"` でメンバーが更新される
- [ ] `./gho members delete <id>` でメンバーが削除される

#### 6. ドキュメント更新

- [ ] `docs/PROJECT_STATUS.md` を更新
- [ ] `docs/NEXT_STEPS.md` を更新（Phase 5に移行）
- [ ] `docs/ARCHITECTURE.md` を更新
- [ ] `README.md` にMembersコマンドの使用例を追加

#### 7. コミット

- [ ] Phase 4完了のコミットを作成

### 実装の開始方法

```bash
# featureブランチを作成
git checkout -b feature/phase4-members

# Members APIのテストから開始（TDD）
# internal/ghostapi/members_test.go を作成
```

### 参考: Ghost Admin API仕様

**Members API**:
- エンドポイント: `/ghost/api/admin/members/`
- メソッド: GET, POST, PUT, DELETE
- パラメータ: `limit`, `filter`, `order`, `include`
- 主要フィールド: `email` (必須), `name`, `note`, `subscribed`, `labels`

詳細: https://ghost.org/docs/admin-api/

### 実装時の注意点

1. **TDD原則を厳守**
   - テストを先に書く（RED）
   - 最小限の実装（GREEN）
   - リファクタリング（REFACTOR）

2. **エラーハンドリング**
   - エラーメッセージは日本語で具体的に
   - エラーのラップに `fmt.Errorf` と `%w` を使用
   - 画像アップロード時のファイルサイズ制限を考慮

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

6. **Members特有の注意点**
   - emailアドレスの検証（RFC 5322準拠）
   - 重複emailチェック（Ghost APIがエラーを返す）
   - subscribed、labels、newslettersなどのフィールド対応

## Phase 5以降の予定

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
