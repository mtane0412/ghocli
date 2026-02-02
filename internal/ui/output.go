/**
 * output.go
 * UI output functionality
 *
 * Provides separation of stdout/stderr:
 * - Data output → stdout
 * - Progress messages → stderr
 * - Error messages → stderr
 */

package ui

import (
	"context"
	"fmt"
	"io"
)

// Output manages UI output
type Output struct {
	stdout io.Writer
	stderr io.Writer
}

// context key type
type contextKey int

const (
	uiKey contextKey = iota
)

// WithUI sets UI output in the context
func WithUI(ctx context.Context, ui *Output) context.Context {
	return context.WithValue(ctx, uiKey, ui)
}

// FromContext retrieves UI output from the context
// Returns nil if UI is not set
func FromContext(ctx context.Context) *Output {
	if ui, ok := ctx.Value(uiKey).(*Output); ok {
		return ui
	}
	return nil
}

// NewOutput creates a new Output
func NewOutput(stdout, stderr io.Writer) *Output {
	return &Output{
		stdout: stdout,
		stderr: stderr,
	}
}

// PrintData outputs data to stdout
func (o *Output) PrintData(data string) error {
	_, err := fmt.Fprintln(o.stdout, data)
	return err
}

// PrintMessage outputs a progress message to stderr
func (o *Output) PrintMessage(message string) {
	fmt.Fprintln(o.stderr, message)
}

// PrintError outputs an error message to stderr
func (o *Output) PrintError(message string) {
	fmt.Fprintln(o.stderr, message)
}
