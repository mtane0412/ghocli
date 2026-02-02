/**
 * confirm_test.go
 * Test code for destructive operation confirmation mechanism
 */

package cmd

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/mtane0412/ghocli/internal/ui"
)

// TestConfirmDestructive_SkipConfirmationWhenForceEnabled verifies that confirmation is skipped when Force flag is enabled
func TestConfirmDestructive_SkipConfirmationWhenForceEnabled(t *testing.T) {
	// Precondition: Prepare RootFlags with Force flag enabled
	root := &RootFlags{
		Force:   true,
		NoInput: false,
	}

	// Embed UI in context
	var stdout, stderr bytes.Buffer
	output := ui.NewOutput(&stdout, &stderr)
	ctx := ui.WithUI(context.Background(), output)

	// Execute: Call ConfirmDestructive
	err := ConfirmDestructive(ctx, root, "delete post 'Test Article'")

	// Verify: No error should occur
	if err != nil {
		t.Errorf("When Force=true, no error is expected, but error was returned: %v", err)
	}
}

// TestConfirmDestructive_FunctionExistsEvenWhenForceDisabled verifies that function exists even when Force flag is disabled
func TestConfirmDestructive_FunctionExistsEvenWhenForceDisabled(t *testing.T) {
	// Precondition: Prepare RootFlags with Force flag disabled
	root := &RootFlags{
		Force:   false,
		NoInput: false,
	}

	// Embed UI in context
	var stdout, stderr bytes.Buffer
	output := ui.NewOutput(&stdout, &stderr)
	ctx := ui.WithUI(context.Background(), output)

	// Execute: At this point, implementation is incomplete, so only verify function can be called
	// Actual interactive input test requires separate implementation
	_ = ConfirmDestructive(ctx, root, "delete post 'Test Article'")
}

// TestConfirmDestructive_ErrorWhenNoInputEnabledWithoutForce verifies that error occurs when NoInput flag is enabled without Force
func TestConfirmDestructive_ErrorWhenNoInputEnabledWithoutForce(t *testing.T) {
	// Precondition: Prepare RootFlags with NoInput flag enabled and Force flag disabled
	root := &RootFlags{
		Force:   false,
		NoInput: true,
	}

	// Embed UI in context
	var stdout, stderr bytes.Buffer
	output := ui.NewOutput(&stdout, &stderr)
	ctx := ui.WithUI(context.Background(), output)

	// Execute: Call ConfirmDestructive
	err := ConfirmDestructive(ctx, root, "delete post 'Test Article'")

	// Verify: Error should occur
	if err == nil {
		t.Error("When NoInput=true and Force=false, an error is expected, but no error was returned")
	}

	// Verify: ExitError should be returned
	var exitErr *ExitError
	if !errors.As(err, &exitErr) {
		t.Errorf("ExitError is expected, but a different error was returned: %T %v", err, err)
	}

	// Verify: Error message should contain "non-interactive"
	if err != nil && !contains(err.Error(), "non-interactive") {
		t.Errorf("Error message does not contain 'non-interactive': %v", err)
	}
}

// TestConfirmDestructive_SkipConfirmationWhenBothNoInputAndForceEnabled verifies that confirmation is skipped when both NoInput and Force are enabled
func TestConfirmDestructive_SkipConfirmationWhenBothNoInputAndForceEnabled(t *testing.T) {
	// Precondition: Prepare RootFlags with both NoInput and Force flags enabled
	root := &RootFlags{
		Force:   true,
		NoInput: true,
	}

	// Embed UI in context
	var stdout, stderr bytes.Buffer
	output := ui.NewOutput(&stdout, &stderr)
	ctx := ui.WithUI(context.Background(), output)

	// Execute: Call ConfirmDestructive
	err := ConfirmDestructive(ctx, root, "delete post 'Test Article'")

	// Verify: No error should occur
	if err != nil {
		t.Errorf("When Force=true, no error is expected, but error was returned: %v", err)
	}
}

// contains is a helper function to check if a substring is contained in a string
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
