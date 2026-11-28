package exitcode

import (
	"errors"
	"fmt"
)

const (
	// Success indicates successful command execution.
	Success = 0

	// General indicates a general runtime error (e.g. server 500, network failure, unexpected error).
	General = 1

	// Usage indicates a command-line usage error (invalid flags, missing required arguments, syntax error).
	Usage = 2

	// Auth indicates an authentication error (unauthenticated, missing token, expired refresh token).
	Auth = 3

	// Forbidden indicates an authorization error (insufficient grants, forbidden operation, attenuated permission mismatch).
	Forbidden = 4

	// NotFound indicates a resource was not found (dataset, table, store, or model does not exist).
	NotFound = 5
)

// ExitError represents an error associated with a specific CLI exit code.
type ExitError struct {
	Code int
	Err  error
}

// Error returns the error message.
func (e *ExitError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("exit code %d", e.Code)
	}
	return e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *ExitError) Unwrap() error {
	return e.Err
}

// New creates an ExitError with the given exit code and formatted message.
func New(code int, message string) error {
	return &ExitError{
		Code: code,
		Err:  errors.New(message),
	}
}

// Wrap encapsulates an existing error with the specified exit code.
func Wrap(code int, err error) error {
	if err == nil {
		return nil
	}
	return &ExitError{
		Code: code,
		Err:  err,
	}
}

// Code extracts the exit code from an error.
// If err is nil, Success (0) is returned.
// If err is an *ExitError, its Code is returned.
// Otherwise, General (1) is returned.
func Code(err error) int {
	if err == nil {
		return Success
	}
	var exitErr *ExitError
	if errors.As(err, &exitErr) {
		return exitErr.Code
	}
	return General
}
