package cmd

import (
	"errors"
	"testing"
)

// TestExitCode verifies the behavior of ExitCode function
func TestExitCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "return 0 when error is nil",
			err:      nil,
			wantCode: 0,
		},
		{
			name:     "return 1 for regular error",
			err:      errors.New("regular error"),
			wantCode: 1,
		},
		{
			name:     "return 2 for ExitError with code 2",
			err:      &ExitError{Code: 2, Err: errors.New("custom error")},
			wantCode: 2,
		},
		{
			name:     "return 130 for ExitError with code 130 (Ctrl+C)",
			err:      &ExitError{Code: 130, Err: errors.New("interrupted")},
			wantCode: 130,
		},
		{
			name:     "return correct code even for wrapped ExitError",
			err:      errors.Join(errors.New("wrapper"), &ExitError{Code: 3, Err: errors.New("inner error")}),
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

// TestExitError_Error verifies the behavior of ExitError.Error() method
func TestExitError_Error(t *testing.T) {
	tests := []struct {
		name        string
		exitErr     *ExitError
		wantMessage string
	}{
		{
			name: "return inner error message",
			exitErr: &ExitError{
				Code: 1,
				Err:  errors.New("test error message"),
			},
			wantMessage: "test error message",
		},
		{
			name: "return complex error message correctly",
			exitErr: &ExitError{
				Code: 2,
				Err:  errors.New("authentication failed: token is invalid"),
			},
			wantMessage: "authentication failed: token is invalid",
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

// TestExitError_Unwrap verifies the behavior of ExitError.Unwrap() method
func TestExitError_Unwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	exitErr := &ExitError{
		Code: 1,
		Err:  innerErr,
	}

	unwrapped := errors.Unwrap(exitErr)
	if unwrapped != innerErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, innerErr)
	}
}
