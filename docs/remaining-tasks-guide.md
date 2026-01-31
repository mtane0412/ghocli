# 残りタスクの実装ガイド

このドキュメントは、ghoとgogcliの設計統一プロジェクトの残りタスク（6, 8, 9）の実装ガイドです。

---

## タスク6: errfmtパッケージの実装

**優先度**: 中
**目的**: ユーザーフレンドリーなエラーメッセージを提供

### 概要

gogcliのerrfmtパッケージは、エラーメッセージを分かりやすくフォーマットし、ユーザーに適切な対処方法を提示します。

### 実装内容

#### 1. 新規ファイルの作成

**ファイル**: `internal/errfmt/errfmt.go`

```go
package errfmt

import (
	"errors"
	"fmt"
)

// Format はエラーメッセージをユーザーフレンドリーな形式にフォーマットする
func Format(err error) string {
	if err == nil {
		return ""
	}

	// 型付きエラーに応じたメッセージを返す
	var authErr *AuthRequiredError
	if errors.As(err, &authErr) {
		return formatAuthRequiredError(authErr)
	}

	// デフォルトはエラーメッセージをそのまま返す
	return err.Error()
}

// AuthRequiredError は認証が必要なことを示すエラー
type AuthRequiredError struct {
	Site string
	Err  error
}

func (e *AuthRequiredError) Error() string {
	return fmt.Sprintf("サイト '%s' の認証が必要です", e.Site)
}

func (e *AuthRequiredError) Unwrap() error {
	return e.Err
}

// formatAuthRequiredError は認証エラーをフォーマットする
func formatAuthRequiredError(err *AuthRequiredError) string {
	return fmt.Sprintf(`%s

認証を追加するには、以下のコマンドを実行してください：
  gho auth add %s
`, err.Error(), err.Site)
}
```

#### 2. テストファイルの作成

**ファイル**: `internal/errfmt/errfmt_test.go`

```go
package errfmt

import (
	"errors"
	"strings"
	"testing"
)

func TestFormat_AuthRequiredError(t *testing.T) {
	err := &AuthRequiredError{
		Site: "https://example.ghost.io",
		Err:  errors.New("token not found"),
	}

	result := Format(err)

	// エラーメッセージが含まれることを確認
	if !strings.Contains(result, "認証が必要です") {
		t.Error("エラーメッセージが含まれていない")
	}

	// ヘルプコマンドが含まれることを確認
	if !strings.Contains(result, "gho auth add") {
		t.Error("ヘルプコマンドが含まれていない")
	}
}

func TestFormat_NilError(t *testing.T) {
	result := Format(nil)
	if result != "" {
		t.Errorf("Format(nil) = %q, want empty string", result)
	}
}

func TestFormat_GenericError(t *testing.T) {
	err := errors.New("generic error")
	result := Format(err)
	if result != "generic error" {
		t.Errorf("Format(generic error) = %q, want %q", result, "generic error")
	}
}
```

#### 3. 既存コードへの統合

各コマンドのエラーハンドリングでerrfmt.Formatを使用するように変更：

```go
// 例: internal/cmd/site.go
func (c *SiteCmd) Run(ctx context.Context, root *RootFlags) error {
	client, err := getAPIClient(root)
	if err != nil {
		// errfmt.Formatでエラーメッセージをフォーマット
		return fmt.Errorf("%s", errfmt.Format(err))
	}
	// ...
}
```

### 拡張可能なエラー型

必要に応じて以下のエラー型を追加できます：

- `CredentialsNotFoundError`: クレデンシャルが見つからない
- `InvalidAPIKeyError`: APIキーの形式が不正
- `NetworkError`: ネットワークエラー（リトライ推奨）
- `RateLimitError`: レート制限エラー（待機時間表示）

---

## タスク8: confirmコマンドのcontext対応

**優先度**: 中
**目的**: ExitErrorを返すように修正し、contextからUIインスタンスを取得

### 概要

現在のconfirmコマンドは直接標準入力を読み取っていますが、contextからUIインスタンスを取得し、ExitErrorを返すように変更します。

### 実装内容

#### 1. confirm.goの修正

**ファイル**: `internal/cmd/confirm.go`

**変更前**:
```go
func ConfirmDestructive(root *RootFlags, message string) error {
	if root.Force {
		return nil
	}
	if root.NoInput {
		return fmt.Errorf("破壊的操作の確認が必要です。--force フラグを使用してください")
	}

	fmt.Fprintf(os.Stderr, "%s (y/N): ", message)
	scanner := bufio.NewScanner(os.Stdin)
	// ...
}
```

**変更後**:
```go
import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/mtane0412/gho/internal/ui"
)

func ConfirmDestructive(ctx context.Context, root *RootFlags, message string) error {
	// Forceフラグが有効な場合は確認をスキップ
	if root.Force {
		return nil
	}

	// NoInputフラグが有効な場合はExitErrorを返す
	if root.NoInput {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("破壊的操作の確認が必要です。--force フラグを使用してください"),
		}
	}

	// contextからUIインスタンスを取得
	output := ui.FromContext(ctx)
	if output == nil {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("UI出力が初期化されていません"),
		}
	}

	// 確認メッセージを表示
	output.PrintMessage(fmt.Sprintf("%s (y/N): ", message))

	// 標準入力から回答を読み取る
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return &ExitError{
			Code: 130, // Ctrl+C
			Err:  fmt.Errorf("確認がキャンセルされました"),
		}
	}

	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if answer != "y" && answer != "yes" {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("操作がキャンセルされました"),
		}
	}

	return nil
}
```

#### 2. テストの更新

**ファイル**: `internal/cmd/confirm_test.go`

contextを渡すように既存のテストを更新：

```go
func TestConfirmDestructive_Forceフラグが有効な場合は確認をスキップ(t *testing.T) {
	ctx := context.Background()
	root := &RootFlags{Force: true}
	err := ConfirmDestructive(ctx, root, "本当に削除しますか？")
	if err != nil {
		t.Errorf("Force=trueの場合、エラーは発生しないはず: %v", err)
	}
}

func TestConfirmDestructive_NoInputフラグが有効な場合はExitErrorを返す(t *testing.T) {
	ctx := context.Background()
	root := &RootFlags{NoInput: true}
	err := ConfirmDestructive(ctx, root, "本当に削除しますか？")

	// ExitErrorが返されることを確認
	var exitErr *ExitError
	if !errors.As(err, &exitErr) {
		t.Error("ExitErrorが返されるべき")
	}

	// 終了コードが1であることを確認
	if exitErr.Code != 1 {
		t.Errorf("終了コード = %d, want 1", exitErr.Code)
	}
}
```

#### 3. 呼び出し側の更新

ConfirmDestructiveを呼び出しているすべてのコマンドで、contextを渡すように変更：

```go
// 例: internal/cmd/posts.go の削除コマンド
func (c *PostsDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// 確認を求める（contextを渡す）
	if err := ConfirmDestructive(ctx, root, "本当に削除しますか？"); err != nil {
		return err
	}
	// ...
}
```

---

## タスク9: inputパッケージの実装

**優先度**: 低
**目的**: 入力抽象化を実装し、テスタビリティを向上

### 概要

現在はconfirm内で直接bufioを使用していますが、inputパッケージとして抽象化することで、テストが容易になります。

### 実装内容

#### 1. 新規パッケージの作成

**ファイル**: `internal/input/prompt.go`

```go
package input

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/ui"
)

// PromptLine はユーザーに入力を促し、1行読み取る
func PromptLine(ctx context.Context, prompt string) (string, error) {
	// contextからUIインスタンスを取得
	output := ui.FromContext(ctx)
	if output != nil {
		output.PrintMessage(prompt)
	} else {
		// UIが設定されていない場合は標準エラー出力に表示
		fmt.Fprint(os.Stderr, prompt)
	}

	// 標準入力から1行読み取る
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("入力の読み取りに失敗: %w", err)
		}
		// Ctrl+C等で中断された場合
		return "", fmt.Errorf("入力がキャンセルされました")
	}

	return scanner.Text(), nil
}

// PromptPassword はパスワード入力を促す（エコーバックなし）
// 注: 実装にはgolang.org/x/term等のライブラリが必要
func PromptPassword(ctx context.Context, prompt string) (string, error) {
	// TODO: 実装（必要に応じて）
	return PromptLine(ctx, prompt)
}
```

#### 2. テストファイルの作成

**ファイル**: `internal/input/prompt_test.go`

```go
package input

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/mtane0412/gho/internal/ui"
)

// 注: 標準入力のモック化が必要なため、このテストは参考実装
func TestPromptLine_基本動作(t *testing.T) {
	t.Skip("標準入力のモック化が必要")

	ctx := context.Background()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := ui.NewOutput(stdout, stderr)
	ctx = ui.WithUI(ctx, output)

	// 標準入力をモックする必要がある
	// （実装方法はテストフレームワークに依存）
}
```

#### 3. confirmコマンドでの使用

**ファイル**: `internal/cmd/confirm.go`

```go
import (
	"context"
	"fmt"
	"strings"

	"github.com/mtane0412/gho/internal/input"
)

func ConfirmDestructive(ctx context.Context, root *RootFlags, message string) error {
	if root.Force {
		return nil
	}
	if root.NoInput {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("破壊的操作の確認が必要です。--force フラグを使用してください"),
		}
	}

	// inputパッケージを使用して入力を取得
	answer, err := input.PromptLine(ctx, fmt.Sprintf("%s (y/N): ", message))
	if err != nil {
		return &ExitError{
			Code: 130,
			Err:  err,
		}
	}

	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer != "y" && answer != "yes" {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("操作がキャンセルされました"),
		}
	}

	return nil
}
```

---

## 実装の順序

以下の順序で実装することを推奨します：

### フェーズ1: errfmtパッケージ（タスク6）
1. `internal/errfmt/errfmt.go` を作成
2. `internal/errfmt/errfmt_test.go` を作成
3. テストを実行して成功を確認
4. 既存コードで使用して効果を確認

### フェーズ2: confirmコマンドの改善（タスク8）
1. `internal/cmd/confirm.go` を修正（contextとExitError対応）
2. `internal/cmd/confirm_test.go` を更新
3. ConfirmDestructiveを呼び出している全コマンドを更新
4. テストを実行して成功を確認

### フェーズ3: inputパッケージ（タスク9）
1. `internal/input/prompt.go` を作成
2. `internal/input/prompt_test.go` を作成（モック化の検討）
3. `internal/cmd/confirm.go` でinputパッケージを使用
4. テストを実行して成功を確認

---

## 品質チェック

各フェーズ完了後、必ず以下を実行してください：

```bash
# ビルド
make build

# テスト
make test

# Lint
make lint

# 型チェック
make type-check
```

すべてクリアしたら、コミットを作成してください。

---

## TDD原則の遵守

このプロジェクトではTDD（テスト駆動開発）を厳格に適用しています：

1. **RED**: 失敗するテストを先に書く
2. **GREEN**: テストを通す最小限のコードを書く
3. **REFACTOR**: コードを整理する

**実装ファーストは禁止**です。必ずテストを先に書いてから実装してください。

---

## 参照リソース

- **gogcliリポジトリ**: 参照実装として利用
- **CLAUDE.md**: プロジェクト固有の開発ルール
- **進捗状況**: `docs/gogcli-alignment-status.md`

---

## 質問・相談

実装中に不明点があれば、以下を参照してください：

1. gogcliの該当パッケージの実装
2. ghoの既存の実装パターン
3. CLAUDE.mdのルール

それでも解決しない場合は、実装を一旦停止して相談してください。
