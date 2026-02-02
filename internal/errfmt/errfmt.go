/**
 * errfmt package
 *
 * Provides functionality to format error messages for users.
 * Generates messages with appropriate solutions for Ghost Admin API specific errors.
 */
package errfmt

import (
	"errors"
	"fmt"
)

// AuthRequiredError represents an error when authentication to Ghost Admin API is required
type AuthRequiredError struct {
	// Site is the Ghost site domain that requires authentication (e.g., "example.ghost.io")
	Site string
	// Err is the original error
	Err error
}

// Error returns the error message for AuthRequiredError
func (e *AuthRequiredError) Error() string {
	if e.Site != "" {
		return fmt.Sprintf("authentication required for %s", e.Site)
	}
	return "authentication required"
}

// Unwrap returns the wrapped original error
func (e *AuthRequiredError) Unwrap() error {
	return e.Err
}

// FormatAuthError formats authentication errors
//
// Generates an error message for sites without configured authentication,
// and suggests using the gho auth add command as a solution.
func FormatAuthError(site string) string {
	return fmt.Sprintf(`No API key configured for site "%s".

Add credentials:
  gho auth add %s https://%s.ghost.io`, site, site, site)
}

// FormatSiteError formats site not specified errors
//
// Generates an error message when no site is specified,
// and suggests using the --site flag or default_site configuration.
func FormatSiteError() string {
	return `No site specified.

Specify with --site flag or set default:
  gho config set default_site myblog`
}

// FormatFlagError formats unknown flag errors
//
// Generates an error message when an unknown flag is specified,
// and suggests using the --help flag.
func FormatFlagError(flag string) string {
	return fmt.Sprintf(`unknown flag %s
Run with --help to see available flags`, flag)
}

// Format formats errors for user-facing output
//
// Recognizes the following special error types and returns messages with appropriate solutions:
// - AuthRequiredError: authentication error → suggests gho auth login command
// - Other errors: returns error message as-is
//
// Returns empty string for nil errors.
func Format(err error) string {
	// Return empty string for nil errors
	if err == nil {
		return ""
	}

	// For AuthRequiredError, suggest authentication method
	var authErr *AuthRequiredError
	if errors.As(err, &authErr) {
		if authErr.Site != "" {
			return fmt.Sprintf(
				"authentication required: %s\n\nSolution:\n  gho auth login %s",
				authErr.Site,
				authErr.Site,
			)
		}
		return "authentication required\n\nSolution:\n  gho auth login <site-url>"
	}

	// Return other errors as-is
	return err.Error()
}
