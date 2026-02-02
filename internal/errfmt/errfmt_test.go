/**
 * errfmt package tests
 *
 * Validates the error formatting functionality.
 */
package errfmt_test

import (
	"errors"
	"testing"

	"github.com/mtane0412/ghocli/internal/errfmt"
)

// TestFormat_NilError verifies that an empty string is returned for nil error
func TestFormat_NilError(t *testing.T) {
	// Execute
	result := errfmt.Format(nil)

	// Verify: empty string is returned
	if result != "" {
		t.Errorf("Format(nil) = %q; want empty string", result)
	}
}

// TestFormat_GenericError verifies that the original message is returned for generic errors
func TestFormat_GenericError(t *testing.T) {
	// Precondition: prepare a generic error
	err := errors.New("something went wrong")

	// Execute
	result := errfmt.Format(err)

	// Verify: error message is returned as-is
	expected := "something went wrong"
	if result != expected {
		t.Errorf("Format(generic error) = %q; want %q", result, expected)
	}
}

// TestFormat_AuthRequiredError verifies that a message with solution is returned for authentication errors
func TestFormat_AuthRequiredError(t *testing.T) {
	// Precondition: prepare an authentication error
	baseErr := errors.New("401 Unauthorized")
	authErr := &errfmt.AuthRequiredError{
		Site: "example.ghost.io",
		Err:  baseErr,
	}

	// Execute
	result := errfmt.Format(authErr)

	// Verify: authentication error message includes solution
	// The message should contain:
	// - Site name
	// - Error content
	// - Solution (gho auth login command)
	if result == "" {
		t.Error("Format(AuthRequiredError) returned empty string")
	}

	// Site name should be included
	if !contains(result, "example.ghost.io") {
		t.Errorf("Format(AuthRequiredError) = %q; want to contain site name 'example.ghost.io'", result)
	}

	// Solution should be included
	if !contains(result, "gho auth login") {
		t.Errorf("Format(AuthRequiredError) = %q; want to contain 'gho auth login'", result)
	}
}

// contains is a helper function that returns whether substr is contained in string s
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

// TestFormatAuthError tests the authentication error format function
func TestFormatAuthError(t *testing.T) {
	// Execute
	msg := errfmt.FormatAuthError("myblog")

	// Verify: site name is included
	if !contains(msg, "myblog") {
		t.Errorf("FormatAuthError() = %q; want to contain 'myblog'", msg)
	}

	// Verify: authentication error message is included
	if !contains(msg, "No API key configured") {
		t.Errorf("FormatAuthError() = %q; want to contain 'No API key configured'", msg)
	}

	// Verify: solution is included
	if !contains(msg, "gho auth add myblog") {
		t.Errorf("FormatAuthError() = %q; want to contain 'gho auth add myblog'", msg)
	}
}

// TestFormatSiteError tests the site not specified error format function
func TestFormatSiteError(t *testing.T) {
	// Execute
	msg := errfmt.FormatSiteError()

	// Verify: site not specified message is included
	if !contains(msg, "No site specified") {
		t.Errorf("FormatSiteError() = %q; want to contain 'No site specified'", msg)
	}

	// Verify: --site flag description is included
	if !contains(msg, "--site") {
		t.Errorf("FormatSiteError() = %q; want to contain '--site'", msg)
	}

	// Verify: config set default_site description is included
	if !contains(msg, "gho config set default_site") {
		t.Errorf("FormatSiteError() = %q; want to contain 'gho config set default_site'", msg)
	}
}

// TestFormatFlagError tests the unknown flag error format function
func TestFormatFlagError(t *testing.T) {
	// Execute
	msg := errfmt.FormatFlagError("--foo")

	// Verify: unknown flag message is included
	if !contains(msg, "unknown flag --foo") {
		t.Errorf("FormatFlagError() = %q; want to contain 'unknown flag --foo'", msg)
	}

	// Verify: --help hint is included
	if !contains(msg, "--help") {
		t.Errorf("FormatFlagError() = %q; want to contain '--help'", msg)
	}
}
