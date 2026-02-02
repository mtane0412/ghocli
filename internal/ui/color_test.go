/**
 * color_test.go
 * Test code for color output functionality
 */

package ui

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestShouldUseColor_NeverMode tests that color output is disabled in Never mode
func TestShouldUseColor_NeverMode(t *testing.T) {
	// Should always return false in Never mode
	assert.False(t, ShouldUseColor(ColorNever), "Should return false in ColorNever mode")
}

// TestShouldUseColor_AlwaysMode tests that color output is enabled in Always mode
func TestShouldUseColor_AlwaysMode(t *testing.T) {
	// Should return true in Always mode when NO_COLOR environment variable is not set
	os.Unsetenv("NO_COLOR")
	assert.True(t, ShouldUseColor(ColorAlways), "Should return true in ColorAlways mode")
}

// TestShouldUseColor_NO_COLOR tests that color output is disabled when NO_COLOR environment variable is set
func TestShouldUseColor_NO_COLOR(t *testing.T) {
	// Set NO_COLOR environment variable
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	// Should return false even in Always mode when NO_COLOR is set
	assert.False(t, ShouldUseColor(ColorAlways), "Should return false even in Always mode when NO_COLOR environment variable is set")
}

// TestShouldUseColor_AutoMode tests that TTY detection is performed in Auto mode
func TestShouldUseColor_AutoMode(t *testing.T) {
	// When NO_COLOR environment variable is not set
	os.Unsetenv("NO_COLOR")

	// In Auto mode, returns the result of TTY detection
	// In CI, it's not a TTY, so it usually returns false
	// This test only verifies that the function can be called
	_ = ShouldUseColor(ColorAuto)
}
