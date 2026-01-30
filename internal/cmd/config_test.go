/**
 * config_test.go
 * configコマンドのテスト
 */

package cmd

import (
	"testing"
)

// TestConfigGetCmd_構造体が存在すること
func TestConfigGetCmd_構造体が存在すること(t *testing.T) {
	var cmd ConfigGetCmd
	_ = cmd
}

// TestConfigSetCmd_構造体が存在すること
func TestConfigSetCmd_構造体が存在すること(t *testing.T) {
	var cmd ConfigSetCmd
	_ = cmd
}

// TestConfigUnsetCmd_構造体が存在すること
func TestConfigUnsetCmd_構造体が存在すること(t *testing.T) {
	var cmd ConfigUnsetCmd
	_ = cmd
}

// TestConfigListCmd_構造体が存在すること
func TestConfigListCmd_構造体が存在すること(t *testing.T) {
	var cmd ConfigListCmd
	_ = cmd
}

// TestConfigPathCmd_構造体が存在すること
func TestConfigPathCmd_構造体が存在すること(t *testing.T) {
	var cmd ConfigPathCmd
	_ = cmd
}

// TestConfigKeysCmd_構造体が存在すること
func TestConfigKeysCmd_構造体が存在すること(t *testing.T) {
	var cmd ConfigKeysCmd
	_ = cmd
}
