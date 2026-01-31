# gho と gogcli の設計統一 - 進捗状況

## 概要

ghoプロジェクトをgogcliの設計パターンに合わせてリファクタリングするプロジェクトの進捗状況を記録します。

## 完了したタスク（高優先度）

### ✅ タスク1: ExitError型の実装
**ファイル**: `internal/cmd/exit.go`, `internal/cmd/exit_test.go`

- 終了コードを管理するExitError型を実装
- ExitCode関数で適切な終了コードを返却
- errors.As対応でエラーチェーンに対応
- テストカバレッジ: 100%

**コミット**: `9f8b4f4` (2026-01-31)

---

### ✅ タスク2: outfmtパッケージのcontext対応
**ファイル**: `internal/outfmt/outfmt.go`, `internal/outfmt/outfmt_test.go`

- Mode構造体（JSON, Plain フラグ）を追加
- WithMode, IsJSON, IsPlain関数でcontextベースの出力モード管理
- tableWriter関数でtabwriterの管理を簡潔化
- 既存のFormatter構造体は互換性のため維持

**コミット**: `9f8b4f4` (2026-01-31)

---

### ✅ タスク3: Execute関数の実装
**ファイル**: `internal/cmd/root.go`, `internal/cmd/root_test.go`

- main.goからロジックを分離したExecute関数を実装
- context初期化、outfmt Mode設定、UI設定を統合
- ExecuteOptionsでバージョン情報を注入可能
- Kongパーサーの構築を一元化

**コミット**: `9f8b4f4` (2026-01-31)

---

### ✅ タスク4: main.goのリファクタリング
**ファイル**: `cmd/gho/main.go`

- Execute関数呼び出しのシンプルなエントリーポイントに変更
- ExitCode関数で適切な終了コードを返却
- buildVersion関数でバージョン情報を構築

**コミット**: `9f8b4f4` (2026-01-31)

---

### ✅ タスク5: 全コマンドのシグネチャ変更
**ファイル**: `internal/cmd/*.go` (17ファイル)

- すべてのRun関数を `Run(ctx context.Context, root *RootFlags) error` に統一
- contextのimportを全ファイルに追加
- キャンセル処理やタイムアウト制御が可能に

**対象ファイル**:
- auth.go, config.go, images.go, members.go, newsletters.go
- offers.go, pages.go, posts.go, site.go, tags.go
- themes.go, tiers.go, users.go, webhooks.go

**コミット**: `9f8b4f4` (2026-01-31)

---

### ✅ タスク7: UIパッケージのcontext対応
**ファイル**: `internal/ui/output.go`, `internal/ui/output_test.go`

- WithUI, FromContext関数でcontextベースのUI管理
- 出力先（stdout/stderr）の分離を維持
- contextから安全にUIインスタンスを取得可能

**コミット**: `9f8b4f4` (2026-01-31)

---

## 残りのタスク（中〜低優先度）

### ⏳ タスク6: errfmtパッケージの実装
**優先度**: 中

**目的**: ユーザーフレンドリーなエラーメッセージを提供

**詳細**: `docs/remaining-tasks-guide.md` の「タスク6」を参照

---

### ⏳ タスク8: confirmコマンドのcontext対応
**優先度**: 中

**目的**: ExitErrorを返すように修正し、contextからUIインスタンスを取得

**詳細**: `docs/remaining-tasks-guide.md` の「タスク8」を参照

---

### ⏳ タスク9: inputパッケージの実装
**優先度**: 低

**目的**: 入力抽象化を実装し、テスタビリティを向上

**詳細**: `docs/remaining-tasks-guide.md` の「タスク9」を参照

---

## 品質確認

すべての変更は以下の品質基準をクリアしています：

- ✅ **ビルド成功**: `make build`
- ✅ **全テスト成功**: `make test`
- ✅ **Lint成功**: `make lint` (0 issues)
- ✅ **型チェック成功**: `make type-check`

---

## 次のステップ

残りのタスクを実装する場合は、以下の順序を推奨します：

1. **タスク6 (errfmt)**: エラーメッセージの改善はユーザー体験に直結
2. **タスク8 (confirm)**: 既存のconfirmコマンドの改善
3. **タスク9 (input)**: 入力抽象化は最後に実装しても問題なし

詳細な実装ガイドは `docs/remaining-tasks-guide.md` を参照してください。

---

## 参照リソース

- **元の設計差異分析レポート**: プランモードのトランスクリプト参照
- **gogcliリポジトリ**: 参照実装として利用
- **実装ガイド**: `docs/remaining-tasks-guide.md`
