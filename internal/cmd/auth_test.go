/**
 * auth_test.go
 * authコマンド拡張のテスト
 */

package cmd

import (
	"testing"
)

// TestAuthTokensCmd_構造体が存在すること
func TestAuthTokensCmd_構造体が存在すること(t *testing.T) {
	var cmd AuthTokensCmd
	_ = cmd
}

// TestAuthTokensListCmd_構造体が存在すること
func TestAuthTokensListCmd_構造体が存在すること(t *testing.T) {
	var cmd AuthTokensListCmd
	_ = cmd
}

// TestAuthTokensDeleteCmd_構造体が存在すること
func TestAuthTokensDeleteCmd_構造体が存在すること(t *testing.T) {
	var cmd AuthTokensDeleteCmd
	_ = cmd
}

// TestAuthCredentialsCmd_構造体が存在すること
func TestAuthCredentialsCmd_構造体が存在すること(t *testing.T) {
	var cmd AuthCredentialsCmd
	_ = cmd
}
