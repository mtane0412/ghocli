/**
 * offers_test.go
 * Test code for offer management commands
 *
 * Includes tests for new commands added in Phase 2.
 */

package cmd

import (
	"testing"
)

// TestOffersInfoCmd_StructExists verifies that OffersInfoCmd struct exists
func TestOffersInfoCmd_StructExists(t *testing.T) {
	// Verify that OffersInfoCmd is defined
	_ = &OffersInfoCmd{}
}

// TestOffersArchiveCmd_StructExists verifies that OffersArchiveCmd struct exists
func TestOffersArchiveCmd_StructExists(t *testing.T) {
	// Verify that OffersArchiveCmd is defined
	_ = &OffersArchiveCmd{}
}
