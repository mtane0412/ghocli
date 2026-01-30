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

## Phase 6: Newsletters/Tiers/Offers

### 目標

Newsletters（ニュースレター）、Tiers（サブスクリプションプラン）、Offers（オファー）の管理機能を実装し、Ghost Admin APIのビジネス機能を強化する

### タスクリスト

#### 1. Newsletters API実装

- [ ] `internal/ghostapi/newsletters.go` を作成
  - [ ] Newsletter型定義（ID、Name、Slug、Statusなど）
  - [ ] NewsletterListOptions型定義
  - [ ] テスト作成（`newsletters_test.go`）
  - [ ] `ListNewsletters(options NewsletterListOptions) (*NewsletterListResponse, error)` 実装
  - [ ] `GetNewsletter(idOrSlug string) (*Newsletter, error)` 実装

#### 2. Tiers API実装

- [ ] `internal/ghostapi/tiers.go` を作成
  - [ ] Tier型定義（ID、Name、Slug、Type、Priceなど）
  - [ ] TierListOptions型定義
  - [ ] テスト作成（`tiers_test.go`）
  - [ ] `ListTiers(options TierListOptions) (*TierListResponse, error)` 実装
  - [ ] `GetTier(idOrSlug string) (*Tier, error)` 実装

#### 3. Offers API実装

- [ ] `internal/ghostapi/offers.go` を作成
  - [ ] Offer型定義（ID、Name、Code、Tier、Discountなど）
  - [ ] OfferListOptions型定義
  - [ ] テスト作成（`offers_test.go`）
  - [ ] `ListOffers(options OfferListOptions) (*OfferListResponse, error)` 実装
  - [ ] `GetOffer(id string) (*Offer, error)` 実装

#### 4. コマンド実装

- [ ] `internal/cmd/newsletters.go` を作成
- [ ] `internal/cmd/tiers.go` を作成
- [ ] `internal/cmd/offers.go` を作成

#### 5. CLIに統合

- [ ] `internal/cmd/root.go` に各コマンドを追加

#### 6. 品質チェック & ドキュメント更新

- [ ] すべてのテストがパス（`make test`）
- [ ] 型チェック成功（`make type-check`）
- [ ] ビルド成功（`make build`）
- [ ] `docs/PROJECT_STATUS.md` を更新
- [ ] Phase 6完了のコミットを作成

### 参考: Ghost Admin API仕様

**Newsletters API**:
- エンドポイント: `/ghost/api/admin/newsletters/`
- メソッド: GET
- 主要フィールド: `id`, `name`, `slug`, `status`, `subscribe_on_signup`

**Tiers API**:
- エンドポイント: `/ghost/api/admin/tiers/`
- メソッド: GET
- 主要フィールド: `id`, `name`, `slug`, `type`, `monthly_price`, `yearly_price`

**Offers API**:
- エンドポイント: `/ghost/api/admin/offers/`
- メソッド: GET
- 主要フィールド: `id`, `name`, `code`, `tier`, `discount_type`, `discount_amount`

詳細: https://ghost.org/docs/admin-api/

## Phase 7以降の予定

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
