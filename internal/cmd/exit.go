package cmd

import "errors"

// ExitError は終了コードを持つエラー型
// コマンドの実行結果に応じた適切な終了コードを返すために使用する
type ExitError struct {
	// Code は終了コード（0は成功、1以上はエラー）
	Code int
	// Err は内部エラー
	Err error
}

// Error はerrorインターフェースを実装する
// 内部エラーのメッセージをそのまま返す
func (e *ExitError) Error() string {
	return e.Err.Error()
}

// Unwrap はエラーチェーン対応のためにUnwrapメソッドを実装する
// errors.As/errors.Isで正しく動作するために必要
func (e *ExitError) Unwrap() error {
	return e.Err
}

// ExitCode はエラーから終了コードを取得する
// - errがnilの場合は0を返す
// - errがExitErrorの場合はそのコードを返す
// - それ以外のエラーの場合は1を返す
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *ExitError
	if errors.As(err, &exitErr) {
		return exitErr.Code
	}
	return 1
}
