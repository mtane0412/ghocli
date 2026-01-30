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

✅ **Phase 5: Users管理** - 完了（2026-01-30）

- Users API（ListUsers、GetUser、UpdateUser）※ Create/Delete非サポート
- Usersコマンド（list、get、update）

✅ **Phase 6: Newsletters/Tiers/Offers** - 完了（2026-01-30）

- Newsletters API（ListNewsletters、GetNewsletter）
- Tiers API（ListTiers、GetTier）
- Offers API（ListOffers、GetOffer）
- Newslettersコマンド（list、get）
- Tiersコマンド（list、get）
- Offersコマンド（list、get）

## Phase 7: Themes/Webhooks

### 目標

Themes（テーマ）とWebhooks（Webhook）の管理機能を実装し、Ghost Admin APIの開発者向け機能を強化する

### タスクリスト

#### 1. Themes API実装

- [ ] `internal/ghostapi/themes.go` を作成
  - [ ] Theme型定義（Name、Package、Active、Templatesなど）
  - [ ] ThemeListOptions型定義
  - [ ] テスト作成（`themes_test.go`）
  - [ ] `ListThemes() ([]Theme, error)` 実装
  - [ ] `UploadTheme(file io.Reader, filename string) (*Theme, error)` 実装
  - [ ] `ActivateTheme(name string) (*Theme, error)` 実装
  - [ ] `DeleteTheme(name string) error` 実装

#### 2. Webhooks API実装

- [ ] `internal/ghostapi/webhooks.go` を作成
  - [ ] Webhook型定義（ID、Event、TargetURL、Secretなど）
  - [ ] WebhookListOptions型定義
  - [ ] テスト作成（`webhooks_test.go`）
  - [ ] `ListWebhooks() ([]Webhook, error)` 実装
  - [ ] `CreateWebhook(webhook *Webhook) (*Webhook, error)` 実装
  - [ ] `UpdateWebhook(id string, webhook *Webhook) (*Webhook, error)` 実装
  - [ ] `DeleteWebhook(id string) error` 実装

#### 3. コマンド実装

- [ ] `internal/cmd/themes.go` を作成
- [ ] `internal/cmd/webhooks.go` を作成

#### 4. CLIに統合

- [ ] `internal/cmd/root.go` に各コマンドを追加

#### 5. 品質チェック & ドキュメント更新

- [ ] すべてのテストがパス（`make test`）
- [ ] 型チェック成功（`make type-check`）
- [ ] ビルド成功（`make build`）
- [ ] `docs/PROJECT_STATUS.md` を更新
- [ ] Phase 7完了のコミットを作成

### 参考: Ghost Admin API仕様

**Themes API**:
- エンドポイント: `/ghost/api/admin/themes/`
- メソッド: GET, POST, PUT, DELETE
- 主要フィールド: `name`, `package`, `active`, `templates`

**Webhooks API**:
- エンドポイント: `/ghost/api/admin/webhooks/`
- メソッド: GET, POST, PUT, DELETE
- 主要フィールド: `id`, `event`, `target_url`, `secret`, `name`, `api_version`

詳細: https://ghost.org/docs/admin-api/

## Phase 8以降の予定

Phase 7完了後に検討します。

## 質問・相談

実装中に疑問が生じた場合：

1. `docs/ARCHITECTURE.md` でアーキテクチャを確認
2. `docs/DEVELOPMENT_GUIDE.md` で開発ガイドを確認
3. Phase 1の実装を参考にする
4. Ghost Admin APIドキュメントを参照

## フィードバック

実装完了後：

1. `docs/PROJECT_STATUS.md` を更新
2. `docs/NEXT_STEPS.md` を更新（Phase 7に移行）
3. 学んだことや改善点を記録
