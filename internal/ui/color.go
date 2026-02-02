/**
 * color.go
 * Color output functionality
 *
 * Provides color output control that supports the --color flag and NO_COLOR environment variable.
 */

package ui

import (
	"os"

	"github.com/muesli/termenv"
)

// ColorMode represents the color output mode
type ColorMode string

const (
	// ColorAuto enables color output only when TTY
	ColorAuto ColorMode = "auto"
	// ColorAlways always enables color output
	ColorAlways ColorMode = "always"
	// ColorNever disables color output
	ColorNever ColorMode = "never"
)

// ShouldUseColor determines whether color output should be used
func ShouldUseColor(mode ColorMode) bool {
	// Always false for Never mode
	if mode == ColorNever {
		return false
	}

	// False if NO_COLOR environment variable is set
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Always true for Always mode
	if mode == ColorAlways {
		return true
	}

	// For Auto mode, determine based on TTY
	profile := termenv.NewOutput(os.Stdout).Profile
	return profile != termenv.Ascii
}
