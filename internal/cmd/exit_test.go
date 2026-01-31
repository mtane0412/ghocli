package cmd

import (
	"errors"
	"testing"
)

// TestExitCode はExitCode関数の動作を検証する
func TestExitCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "エラーがnilの場合は0を返す",
			err:      nil,
			wantCode: 0,
		},
		{
			name:     "通常のエラーの場合は1を返す",
			err:      errors.New("通常のエラー"),
			wantCode: 1,
		},
		{
			name:     "ExitErrorでコード2の場合は2を返す",
			err:      &ExitError{Code: 2, Err: errors.New("カスタムエラー")},
			wantCode: 2,
		},
		{
			name:     "ExitErrorでコード130の場合は130を返す（Ctrl+C）",
			err:      &ExitError{Code: 130, Err: errors.New("中断")},
			wantCode: 130,
		},
		{
			name:     "ラップされたExitErrorでも正しくコードを返す",
			err:      errors.Join(errors.New("ラッパー"), &ExitError{Code: 3, Err: errors.New("内部エラー")}),
			wantCode: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode := ExitCode(tt.err)
			if gotCode != tt.wantCode {
				t.Errorf("ExitCode() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

// TestExitError_Error はExitError.Error()メソッドの動作を検証する
func TestExitError_Error(t *testing.T) {
	tests := []struct {
		name        string
		exitErr     *ExitError
		wantMessage string
	}{
		{
			name: "内部エラーのメッセージを返す",
			exitErr: &ExitError{
				Code: 1,
				Err:  errors.New("テストエラーメッセージ"),
			},
			wantMessage: "テストエラーメッセージ",
		},
		{
			name: "複雑なエラーメッセージも正しく返す",
			exitErr: &ExitError{
				Code: 2,
				Err:  errors.New("認証に失敗しました: トークンが無効です"),
			},
			wantMessage: "認証に失敗しました: トークンが無効です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessage := tt.exitErr.Error()
			if gotMessage != tt.wantMessage {
				t.Errorf("ExitError.Error() = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}

// TestExitError_Unwrap はExitError.Unwrap()メソッドの動作を検証する
func TestExitError_Unwrap(t *testing.T) {
	innerErr := errors.New("内部エラー")
	exitErr := &ExitError{
		Code: 1,
		Err:  innerErr,
	}

	unwrapped := errors.Unwrap(exitErr)
	if unwrapped != innerErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, innerErr)
	}
}
