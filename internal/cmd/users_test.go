/**
 * users_test.go
 * Test code for user management commands
 */

package cmd

import (
	"testing"
)

// TestUsersInfoCmd_StructExists verifies that UsersInfoCmd struct exists
func TestUsersInfoCmd_StructExists(t *testing.T) {
	// Verify that UsersInfoCmd is defined
	_ = &UsersInfoCmd{}
}
