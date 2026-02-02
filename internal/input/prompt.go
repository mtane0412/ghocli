/**
 * prompt.go
 * User input prompt functionality
 *
 * Provides functionality to interactively prompt users for input on the command line.
 */
package input

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mtane0412/gho/internal/ui"
)

// PromptLine reads a single line of input from os.Stdin
//
// ctx: Context (used to retrieve UI)
// prompt: Prompt string to display to the user
//
// Returns:
//   - The read line (without line ending characters)
//   - An error if reading fails
func PromptLine(ctx context.Context, prompt string) (string, error) {
	return PromptLineFrom(ctx, prompt, os.Stdin)
}

// PromptLineFrom reads a single line of input from the specified io.Reader
//
// ctx: Context (used to retrieve UI)
// prompt: Prompt string to display to the user
// r: The input io.Reader
//
// Returns:
//   - The read line (without line ending characters)
//   - An error if reading fails
func PromptLineFrom(ctx context.Context, prompt string, r io.Reader) (string, error) {
	// Get UI from context and output the prompt
	if u := ui.FromContext(ctx); u != nil {
		// If UI is set, use PrintMessage
		u.PrintMessage(prompt)
	} else {
		// If UI is not set, output directly to stderr
		_, _ = fmt.Fprint(os.Stderr, prompt)
	}

	// Read a single line
	return ReadLine(r)
}

// PromptPassword displays a prompt for password input (TODO: implementation pending)
//
// This function is planned for future implementation. Currently unimplemented.
func PromptPassword(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement password input using terminal.ReadPassword
	return "", fmt.Errorf("PromptPassword is not implemented yet")
}
