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

- Newsletters API（ListNewsletters、GetNewsletter、CreateNewsletter、UpdateNewsletter）
- Tiers API（ListTiers、GetTier、CreateTier、UpdateTier）
- Offers API（ListOffers、GetOffer、CreateOffer、UpdateOffer）
- Newslettersコマンド（list、get、create、update）
- Tiersコマンド（list、get、create、update）
- Offersコマンド（list、get、create、update）
- 破壊的操作の確認機構（`--force`フラグでスキップ可能）

✅ **Phase 7: Themes/Webhooks** - 完了（2026-01-30）

- Themes API（ListThemes、UploadTheme、ActivateTheme）
- Webhooks API（CreateWebhook、UpdateWebhook、DeleteWebhook）※ List/Get非サポート
- Themesコマンド（list、upload、activate）
- Webhooksコマンド（create、update、delete）

## Phase 8以降の予定

現時点で主要なGhost Admin API機能の実装が完了しました。今後、以下の拡張機能を検討できます：

### 考えられる拡張機能

1. **データエクスポート/インポート機能**
   - コンテンツのバックアップ/リストア機能
   - 他のブログプラットフォームからの移行支援

2. **バッチ操作機能**
   - 複数の投稿/ページの一括更新
   - タグの一括割り当て
   - メンバーの一括インポート

3. **検索・フィルタリングの拡張**
   - 高度な検索クエリビルダー
   - カスタムフィルタのプリセット保存

4. **レポート機能**
   - サイト統計の表示
   - メンバーレポート
   - コンテンツレポート

5. **対話的UIモード**
   - インタラクティブな投稿エディタ
   - TUIベースのブラウザ

6. **CI/CD統合**
   - GitHub Actionsワークフロー例
   - 自動デプロイスクリプト

### 次のアクション

実装優先度や必要性に応じて、上記の拡張機能から選択するか、ユーザーフィードバックに基づいて新機能を検討します。

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
