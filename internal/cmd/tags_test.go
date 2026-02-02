/**
 * tags_test.go
 * Test code for tag management commands
 */

package cmd

import (
	"testing"
)

// TestTagsInfoCmd_StructExists verifies that TagsInfoCmd struct exists
func TestTagsInfoCmd_StructExists(t *testing.T) {
	// Verify that TagsInfoCmd is defined
	_ = &TagsInfoCmd{}
}
