/**
 * confirm.go
 * Destructive operation confirmation mechanism
 *
 * Based on gogcli's safety mechanism, performs user confirmation before executing destructive operations
 */

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/mtane0412/ghocli/internal/input"
	"github.com/mtane0412/ghocli/internal/ui"
)

// ConfirmDestructive performs user confirmation before executing destructive operations
//
// ctx: Context (used to retrieve UI)
// root: RootFlags to retrieve Force and NoInput flags
// message: Description of the operation to execute (e.g., "delete post 'Test Article'")
//
// Return value:
//   - Always returns nil if Force=true
//   - Returns ExitError{Code: 1} if NoInput=true and Force=false
//   - In non-interactive environments (not a TTY), returns ExitError{Code: 1} without Force
//   - In interactive environments, displays a confirmation prompt and returns ExitError{Code: 1} for inputs other than y/yes
func ConfirmDestructive(ctx context.Context, root *RootFlags, message string) error {
	// Skip confirmation if Force flag is enabled
	if root.Force {
		return nil
	}

	// If NoInput flag is enabled or in non-interactive environment, prohibit interactive input
	if root.NoInput || !term.IsTerminal(int(os.Stdin.Fd())) {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("refusing to %s without --force (non-interactive)", message),
		}
	}

	// Retrieve UI from context
	output := ui.FromContext(ctx)
	if output == nil {
		// Error if UI is not configured (normally does not occur)
		return &ExitError{
			Code: 1,
			Err:  errors.New("ui not configured in context"),
		}
	}

	// Display confirmation prompt and read user input
	prompt := fmt.Sprintf("Proceed to %s? [y/N]: ", message)
	line, readErr := input.PromptLineFrom(ctx, prompt, os.Stdin)
	if readErr != nil && !errors.Is(readErr, os.ErrClosed) {
		// Treat EOF as cancellation
		if errors.Is(readErr, io.EOF) {
			return &ExitError{Code: 1, Err: errors.New("cancelled")}
		}
		// Other errors
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("failed to read user input: %w", readErr),
		}
	}

	// Normalize input (remove leading/trailing whitespace, convert to lowercase)
	ans := strings.ToLower(strings.TrimSpace(line))

	// Continue only if "y" or "yes"
	if ans == "y" || ans == "yes" {
		return nil
	}

	// Cancel otherwise
	return &ExitError{Code: 1, Err: errors.New("cancelled")}
}
