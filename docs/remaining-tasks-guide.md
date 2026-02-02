# Remaining Tasks Implementation Guide

This document provides implementation guidance for the remaining tasks (6, 8, 9) of the gho and gogcli design alignment project.

---

## Task 6: errfmt Package Implementation

**Priority**: Medium
**Purpose**: Provide user-friendly error messages

### Overview

The gogcli errfmt package formats error messages to be clear and suggests appropriate solutions to users.

### Implementation

#### 1. Create New Files

**File**: `internal/errfmt/errfmt.go`

```go
package errfmt

import (
	"errors"
	"fmt"
)

// Format formats error messages into user-friendly format
func Format(err error) string {
	if err == nil {
		return ""
	}

	// Return messages according to typed errors
	var authErr *AuthRequiredError
	if errors.As(err, &authErr) {
		return formatAuthRequiredError(authErr)
	}

	// Default returns error message as-is
	return err.Error()
}

// AuthRequiredError indicates authentication is required
type AuthRequiredError struct {
	Site string
	Err  error
}

func (e *AuthRequiredError) Error() string {
	return fmt.Sprintf("authentication required for site '%s'", e.Site)
}

func (e *AuthRequiredError) Unwrap() error {
	return e.Err
}

// formatAuthRequiredError formats authentication errors
func formatAuthRequiredError(err *AuthRequiredError) string {
	return fmt.Sprintf(`%s

To add authentication, run:
  gho auth add %s
`, err.Error(), err.Site)
}
```

#### 2. Create Test File

**File**: `internal/errfmt/errfmt_test.go`

```go
package errfmt

import (
	"errors"
	"strings"
	"testing"
)

func TestFormat_AuthRequiredError(t *testing.T) {
	err := &AuthRequiredError{
		Site: "https://example.ghost.io",
		Err:  errors.New("token not found"),
	}

	result := Format(err)

	// Verify error message is included
	if !strings.Contains(result, "authentication required") {
		t.Error("Error message not included")
	}

	// Verify help command is included
	if !strings.Contains(result, "gho auth add") {
		t.Error("Help command not included")
	}
}

func TestFormat_NilError(t *testing.T) {
	result := Format(nil)
	if result != "" {
		t.Errorf("Format(nil) = %q, want empty string", result)
	}
}

func TestFormat_GenericError(t *testing.T) {
	err := errors.New("generic error")
	result := Format(err)
	if result != "generic error" {
		t.Errorf("Format(generic error) = %q, want %q", result, "generic error")
	}
}
```

#### 3. Integration into Existing Code

Modify error handling in each command to use errfmt.Format:

```go
// Example: internal/cmd/site.go
func (c *SiteCmd) Run(ctx context.Context, root *RootFlags) error {
	client, err := getAPIClient(root)
	if err != nil {
		// Format error message with errfmt.Format
		return fmt.Errorf("%s", errfmt.Format(err))
	}
	// ...
}
```

### Extensible Error Types

Add the following error types as needed:

- `CredentialsNotFoundError`: Credentials not found
- `InvalidAPIKeyError`: Invalid API key format
- `NetworkError`: Network error (retry recommended)
- `RateLimitError`: Rate limit error (display wait time)

---

## Task 8: Context Support for confirm Command

**Priority**: Medium
**Purpose**: Modify to return ExitError and retrieve UI instance from context

### Overview

Currently, the confirm command reads directly from standard input, but we'll change it to retrieve UI instance from context and return ExitError.

### Implementation

#### 1. Modify confirm.go

**File**: `internal/cmd/confirm.go`

**Before**:
```go
func ConfirmDestructive(root *RootFlags, message string) error {
	if root.Force {
		return nil
	}
	if root.NoInput {
		return fmt.Errorf("confirmation required for destructive operation. Use --force flag")
	}

	fmt.Fprintf(os.Stderr, "%s (y/N): ", message)
	scanner := bufio.NewScanner(os.Stdin)
	// ...
}
```

**After**:
```go
import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/mtane0412/gho/internal/ui"
)

func ConfirmDestructive(ctx context.Context, root *RootFlags, message string) error {
	// Skip confirmation if Force flag is enabled
	if root.Force {
		return nil
	}

	// Return ExitError if NoInput flag is enabled
	if root.NoInput {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("confirmation required for destructive operation. Use --force flag"),
		}
	}

	// Retrieve UI instance from context
	output := ui.FromContext(ctx)
	if output == nil {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("UI output not initialized"),
		}
	}

	// Display confirmation message
	output.PrintMessage(fmt.Sprintf("%s (y/N): ", message))

	// Read response from standard input
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return &ExitError{
			Code: 130, // Ctrl+C
			Err:  fmt.Errorf("confirmation cancelled"),
		}
	}

	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if answer != "y" && answer != "yes" {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("operation cancelled"),
		}
	}

	return nil
}
```

#### 2. Update Tests

**File**: `internal/cmd/confirm_test.go`

Update existing tests to pass context:

```go
func TestConfirmDestructive_SkipsConfirmationWhenForceEnabled(t *testing.T) {
	ctx := context.Background()
	root := &RootFlags{Force: true}
	err := ConfirmDestructive(ctx, root, "Really delete?")
	if err != nil {
		t.Errorf("Should not error when Force=true: %v", err)
	}
}

func TestConfirmDestructive_ReturnsExitErrorWhenNoInputEnabled(t *testing.T) {
	ctx := context.Background()
	root := &RootFlags{NoInput: true}
	err := ConfirmDestructive(ctx, root, "Really delete?")

	// Verify ExitError is returned
	var exitErr *ExitError
	if !errors.As(err, &exitErr) {
		t.Error("Should return ExitError")
	}

	// Verify exit code is 1
	if exitErr.Code != 1 {
		t.Errorf("Exit code = %d, want 1", exitErr.Code)
	}
}
```

#### 3. Update Callers

Update all commands calling ConfirmDestructive to pass context:

```go
// Example: internal/cmd/posts.go delete command
func (c *PostsDeleteCmd) Run(ctx context.Context, root *RootFlags) error {
	// Request confirmation (pass context)
	if err := ConfirmDestructive(ctx, root, "Really delete?"); err != nil {
		return err
	}
	// ...
}
```

---

## Task 9: input Package Implementation

**Priority**: Low
**Purpose**: Implement input abstraction to improve testability

### Overview

Currently using bufio directly in confirm, but abstracting as input package makes testing easier.

### Implementation

#### 1. Create New Package

**File**: `internal/input/prompt.go`

```go
package input

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/mtane0412/gho/internal/ui"
)

// PromptLine prompts user and reads one line
func PromptLine(ctx context.Context, prompt string) (string, error) {
	// Retrieve UI instance from context
	output := ui.FromContext(ctx)
	if output != nil {
		output.PrintMessage(prompt)
	} else {
		// Display on stderr if UI not configured
		fmt.Fprint(os.Stderr, prompt)
	}

	// Read one line from standard input
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		// Interrupted by Ctrl+C etc.
		return "", fmt.Errorf("input cancelled")
	}

	return scanner.Text(), nil
}

// PromptPassword prompts for password input (no echo)
// Note: Requires library like golang.org/x/term for implementation
func PromptPassword(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement (if needed)
	return PromptLine(ctx, prompt)
}
```

#### 2. Create Test File

**File**: `internal/input/prompt_test.go`

```go
package input

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/mtane0412/gho/internal/ui"
)

// Note: This test is reference implementation as stdin mocking is required
func TestPromptLine_BasicOperation(t *testing.T) {
	t.Skip("stdin mocking required")

	ctx := context.Background()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := ui.NewOutput(stdout, stderr)
	ctx = ui.WithUI(ctx, output)

	// Need to mock standard input
	// (implementation depends on test framework)
}
```

#### 3. Use in confirm Command

**File**: `internal/cmd/confirm.go`

```go
import (
	"context"
	"fmt"
	"strings"

	"github.com/mtane0412/gho/internal/input"
)

func ConfirmDestructive(ctx context.Context, root *RootFlags, message string) error {
	if root.Force {
		return nil
	}
	if root.NoInput {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("confirmation required for destructive operation. Use --force flag"),
		}
	}

	// Use input package to get input
	answer, err := input.PromptLine(ctx, fmt.Sprintf("%s (y/N): ", message))
	if err != nil {
		return &ExitError{
			Code: 130,
			Err:  err,
		}
	}

	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer != "y" && answer != "yes" {
		return &ExitError{
			Code: 1,
			Err:  fmt.Errorf("operation cancelled"),
		}
	}

	return nil
}
```

---

## Implementation Order

Recommended implementation order:

### Phase 1: errfmt Package (Task 6)
1. Create `internal/errfmt/errfmt.go`
2. Create `internal/errfmt/errfmt_test.go`
3. Run tests to verify success
4. Use in existing code to verify effectiveness

### Phase 2: confirm Command Improvements (Task 8)
1. Modify `internal/cmd/confirm.go` (context and ExitError support)
2. Update `internal/cmd/confirm_test.go`
3. Update all commands calling ConfirmDestructive
4. Run tests to verify success

### Phase 3: input Package (Task 9)
1. Create `internal/input/prompt.go`
2. Create `internal/input/prompt_test.go` (consider mocking)
3. Use input package in `internal/cmd/confirm.go`
4. Run tests to verify success

---

## Quality Checks

After completing each phase, always execute:

```bash
# Build
make build

# Tests
make test

# Lint
make lint

# Type check
make type-check
```

Create commit after all pass.

---

## TDD Principles Adherence

This project strictly applies TDD (Test-Driven Development):

1. **RED**: Write a failing test first
2. **GREEN**: Write minimal code to make the test pass
3. **REFACTOR**: Clean up the code

**Implementation-first is prohibited**. Always write tests first, then implement.

---

## Reference Resources

- **gogcli Repository**: Used as reference implementation
- **CLAUDE.md**: Project-specific development rules
- **Progress Status**: `docs/gogcli-alignment-status.md`

---

## Questions & Consultation

If unclear during implementation, refer to:

1. gogcli implementation of corresponding package
2. gho's existing implementation patterns
3. CLAUDE.md rules

If still unresolved, pause implementation and consult.
