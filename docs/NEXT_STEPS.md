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

## Phase 3: タクソノミー + メディア

### 目標

Tags APIとImages APIを実装し、Ghost Admin APIの基本的なコンテンツ管理機能を完成させる

### タスクリスト

#### 1. Tags API実装

- [ ] `internal/ghostapi/tags.go` を作成
  - [ ] Tag型定義（ID、Name、Slug、Description、VisibilityなどGhost Admin APIのTag構造に準拠）
  - [ ] TagListOptions型定義（Limit、Filter、Orderなど）
  - [ ] テスト作成（`tags_test.go`）
  - [ ] `ListTags(options TagListOptions) ([]Tag, error)` 実装
  - [ ] `GetTag(idOrSlug string) (*Tag, error)` 実装
  - [ ] `CreateTag(tag *Tag) (*Tag, error)` 実装
  - [ ] `UpdateTag(id string, tag *Tag) (*Tag, error)` 実装
  - [ ] `DeleteTag(id string) error` 実装

#### 2. Tagsコマンド実装

- [ ] `internal/cmd/tags.go` を作成
  - [ ] TagsCmd構造体定義
  - [ ] TagsListCmd実装
    - [ ] `--limit` フラグ
    - [ ] `--filter` フラグ（オプション）
  - [ ] TagsGetCmd実装
  - [ ] TagsCreateCmd実装
    - [ ] `--name` フラグ（必須）
    - [ ] `--slug` フラグ（オプション、自動生成）
    - [ ] `--description` フラグ（オプション）
  - [ ] TagsUpdateCmd実装
  - [ ] TagsDeleteCmd実装

#### 3. Images API実装

- [ ] `internal/ghostapi/images.go` を作成
  - [ ] Image型定義（URL、Refなど）
  - [ ] テスト作成（`images_test.go`）
  - [ ] `UploadImage(filePath string) (*Image, error)` 実装
    - [ ] ファイル読み込み
    - [ ] multipart/form-dataでアップロード
    - [ ] アップロード後のURL取得

#### 4. Imagesコマンド実装

- [ ] `internal/cmd/images.go` を作成
  - [ ] ImagesCmd構造体定義
  - [ ] ImagesUploadCmd実装
    - [ ] ファイルパス引数（必須）
    - [ ] ファイル存在確認
    - [ ] 画像形式検証（jpg、png、gif、webpなど）

#### 5. CLIに統合

- [ ] `internal/cmd/root.go` に TagsCmd と ImagesCmd を追加

#### 6. 品質チェック

- [ ] すべてのテストがパス（`make test`）
- [ ] 型チェック成功（`make type-check`）
- [ ] Lint成功（`make lint`）
- [ ] ビルド成功（`make build`）

#### 7. 動作確認

- [ ] `./gho tags list` でTags一覧が表示される
- [ ] `./gho tags get <slug>` でタグ詳細が表示される
- [ ] `./gho tags create --name "Technology"` でタグが作成される
- [ ] `./gho tags update <id> --name "Tech"` でタグが更新される
- [ ] `./gho tags delete <id>` でタグが削除される
- [ ] `./gho images upload path/to/image.jpg` で画像がアップロードされる
- [ ] アップロード後にURLが表示される

#### 8. ドキュメント更新

- [ ] `docs/PROJECT_STATUS.md` を更新
- [ ] `docs/NEXT_STEPS.md` を更新（Phase 4に移行）
- [ ] `README.md` にTags/Imagesコマンドの使用例を追加

#### 9. コミット

- [ ] Phase 3完了のコミットを作成

### 実装の開始方法

```bash
# featureブランチを作成
git checkout -b feature/phase3-taxonomy-media

# Tags APIのテストから開始（TDD）
# internal/ghostapi/tags_test.go を作成
```

### 参考: Ghost Admin API仕様

**Tags API**:
- エンドポイント: `/ghost/api/admin/tags/`
- メソッド: GET, POST, PUT, DELETE
- パラメータ: `limit`, `filter`, `order`, `include`

**Images API**:
- エンドポイント: `/ghost/api/admin/images/upload/`
- メソッド: POST (multipart/form-data)
- フィールド: `file` (画像ファイル)

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

6. **画像アップロード特有の注意点**
   - ファイルサイズ制限（Ghost Admin APIの制限を確認）
   - MIME type検証
   - アップロード進捗表示（大きいファイルの場合）

## Phase 4以降の予定

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
