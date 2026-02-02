/**
 * config_test.go
 * Test code for config command
 */

package cmd

import (
	"testing"
)

// TestConfigGetCmd_StructExists verifies that ConfigGetCmd struct exists
func TestConfigGetCmd_StructExists(t *testing.T) {
	var cmd ConfigGetCmd
	_ = cmd
}

// TestConfigSetCmd_StructExists verifies that ConfigSetCmd struct exists
func TestConfigSetCmd_StructExists(t *testing.T) {
	var cmd ConfigSetCmd
	_ = cmd
}

// TestConfigUnsetCmd_StructExists verifies that ConfigUnsetCmd struct exists
func TestConfigUnsetCmd_StructExists(t *testing.T) {
	var cmd ConfigUnsetCmd
	_ = cmd
}

// TestConfigListCmd_StructExists verifies that ConfigListCmd struct exists
func TestConfigListCmd_StructExists(t *testing.T) {
	var cmd ConfigListCmd
	_ = cmd
}

// TestConfigPathCmd_StructExists verifies that ConfigPathCmd struct exists
func TestConfigPathCmd_StructExists(t *testing.T) {
	var cmd ConfigPathCmd
	_ = cmd
}

// TestConfigKeysCmd_StructExists verifies that ConfigKeysCmd struct exists
func TestConfigKeysCmd_StructExists(t *testing.T) {
	var cmd ConfigKeysCmd
	_ = cmd
}
