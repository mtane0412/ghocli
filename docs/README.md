# gho ドキュメント

## 概要

このディレクトリにはgho（Ghost Admin API CLI）プロジェクトのドキュメントが含まれています。

## ドキュメント一覧

### 📊 [PROJECT_STATUS.md](./PROJECT_STATUS.md)

**目的**: プロジェクトの現在の状態を把握する

**内容**:
- 実装フェーズの進捗状況
- 完了した機能
- 現在のプロジェクト構造
- テストカバレッジ
- 依存関係

**いつ読むか**:
- プロジェクトの状態を確認したいとき
- どこまで実装が完了しているか知りたいとき

---

### 📋 [IMPLEMENTATION_PLAN.md](./IMPLEMENTATION_PLAN.md)

**目的**: 実装計画の全体像を理解する

**内容**:
- 技術スタック
- 全フェーズの実装計画
- 各フェーズの目標と実装内容
- 検証方法
- 開発ワークフロー

**いつ読むか**:
- 次に何を実装すべきか知りたいとき
- 実装の全体像を把握したいとき
- 新しいフェーズを開始する前

---

### 🏗️ [ARCHITECTURE.md](./ARCHITECTURE.md)

**目的**: システムのアーキテクチャを理解する

**内容**:
- プロジェクト構造
- レイヤー構成
- コンポーネント設計
- 認証フロー
- APIリクエストフロー
- エラーハンドリング
- テスト戦略
- セキュリティ考慮事項

**いつ読むか**:
- コードの設計を理解したいとき
- 新しい機能の実装場所を決めるとき
- コードレビュー時

---

### 👨‍💻 [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md)

**目的**: 開発方法を学ぶ

**内容**:
- 開発環境のセットアップ
- 開発ワークフロー（TDD）
- コーディング規約
- テストの書き方
- 品質チェック方法
- 新しいAPIリソースの追加方法
- デバッグ方法
- トラブルシューティング

**いつ読むか**:
- プロジェクトに初めて参加するとき
- コードを書く前
- テストの書き方を確認したいとき
- 品質チェックを実行する前

---

### 🚀 [NEXT_STEPS.md](./NEXT_STEPS.md)

**目的**: 次に何をすべきか確認する

**内容**:
- 現在の状態
- 次のフェーズのタスクリスト
- 実装の開始方法
- 参考情報
- 実装時の注意点

**いつ読むか**:
- 次のタスクを確認したいとき
- 新しいフェーズを開始するとき
- 何から始めればいいか分からないとき

---

## ドキュメントの読み方

### 初めてプロジェクトに参加する場合

1. **PROJECT_STATUS.md** を読んでプロジェクトの状態を把握
2. **ARCHITECTURE.md** を読んでシステム設計を理解
3. **DEVELOPMENT_GUIDE.md** を読んで開発方法を学ぶ
4. **NEXT_STEPS.md** を読んで次のタスクを確認

### 新しいフェーズを開始する場合

1. **NEXT_STEPS.md** でタスクリストを確認
2. **IMPLEMENTATION_PLAN.md** で詳細な計画を確認
3. **DEVELOPMENT_GUIDE.md** で実装方法を確認
4. 実装開始

### コードレビュー時

1. **ARCHITECTURE.md** で設計方針を確認
2. **DEVELOPMENT_GUIDE.md** でコーディング規約を確認

### トラブルシューティング時

1. **DEVELOPMENT_GUIDE.md** のトラブルシューティングセクションを確認
2. **ARCHITECTURE.md** でシステム構造を確認

## ドキュメントの更新

ドキュメントは常に最新の状態に保つ必要があります。

### 更新が必要なタイミング

| ドキュメント | 更新タイミング |
|-------------|---------------|
| PROJECT_STATUS.md | フェーズ完了時 |
| IMPLEMENTATION_PLAN.md | 計画変更時 |
| ARCHITECTURE.md | アーキテクチャ変更時 |
| DEVELOPMENT_GUIDE.md | 開発方法変更時 |
| NEXT_STEPS.md | フェーズ完了時、タスク完了時 |

### 更新手順

1. ドキュメントを編集
2. コミットメッセージに「docs:」プレフィックスを付ける

```bash
git commit -m "docs: PROJECT_STATUS.mdを更新（Phase 2完了）"
```

## フィードバック

ドキュメントの改善提案があれば、以下の方法で共有してください：

1. GitHubのIssueを作成
2. Pull Requestを送信
3. コミットメッセージに記載

## クイックリファレンス

### プロジェクト情報

- **プロジェクト名**: gho
- **説明**: Ghost Admin API CLI
- **言語**: Go 1.22+
- **CLI フレームワーク**: Kong

### ディレクトリ構造

```
gho/
├── cmd/gho/          # エントリーポイント
├── internal/         # 内部パッケージ
│   ├── cmd/         # CLIコマンド
│   ├── config/      # 設定管理
│   ├── secrets/     # キーリング統合
│   ├── ghostapi/    # Ghost API
│   └── outfmt/      # 出力フォーマット
└── docs/            # ドキュメント
```

### 品質チェックコマンド

```bash
make test         # テスト実行
make type-check   # 型チェック
make lint         # Lint実行
make build        # ビルド
```

### 重要なリンク

- [Ghost Admin API Documentation](https://ghost.org/docs/admin-api/)
- [Kong CLI Framework](https://github.com/alecthomas/kong)
- [99designs/keyring](https://github.com/99designs/keyring)
