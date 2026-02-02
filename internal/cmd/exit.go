package cmd

import "errors"

// ExitError is an error type that holds an exit code
// Used to return an appropriate exit code based on command execution results
type ExitError struct {
	// Code is the exit code (0 for success, 1 or higher for errors)
	Code int
	// Err is the internal error
	Err error
}

// Error implements the error interface
// Returns the internal error message as-is
func (e *ExitError) Error() string {
	return e.Err.Error()
}

// Unwrap implements the Unwrap method for error chain support
// Necessary for errors.As/errors.Is to work correctly
func (e *ExitError) Unwrap() error {
	return e.Err
}

// ExitCode retrieves the exit code from an error
// - Returns 0 if err is nil
// - Returns the code if err is ExitError
// - Returns 1 for other errors
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *ExitError
	if errors.As(err, &exitErr) {
		return exitErr.Code
	}
	return 1
}
