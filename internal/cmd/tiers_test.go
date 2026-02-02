/**
 * tiers_test.go
 * Test code for tier management commands
 */

package cmd

import (
	"testing"
)

// TestTiersInfoCmd_StructExists verifies that TiersInfoCmd struct exists
func TestTiersInfoCmd_StructExists(t *testing.T) {
	// Verify that TiersInfoCmd is defined
	_ = &TiersInfoCmd{}
}
