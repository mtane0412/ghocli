/**
 * auth_test.go
 * Test code for auth command extensions
 */

package cmd

import (
	"testing"
)

// TestAuthTokensCmd_StructExists verifies that AuthTokensCmd struct exists
func TestAuthTokensCmd_StructExists(t *testing.T) {
	var cmd AuthTokensCmd
	_ = cmd
}

// TestAuthTokensListCmd_StructExists verifies that AuthTokensListCmd struct exists
func TestAuthTokensListCmd_StructExists(t *testing.T) {
	var cmd AuthTokensListCmd
	_ = cmd
}

// TestAuthTokensDeleteCmd_StructExists verifies that AuthTokensDeleteCmd struct exists
func TestAuthTokensDeleteCmd_StructExists(t *testing.T) {
	var cmd AuthTokensDeleteCmd
	_ = cmd
}

// TestAuthCredentialsCmd_StructExists verifies that AuthCredentialsCmd struct exists
func TestAuthCredentialsCmd_StructExists(t *testing.T) {
	var cmd AuthCredentialsCmd
	_ = cmd
}
