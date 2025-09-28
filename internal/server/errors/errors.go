// Package postgres defines custom error types used throughout the GophKeeper application.
// These errors provide specific error messages for common application scenarios.
package postgres

import "errors"

// ErrUserIsExist indicates that a user with the specified login already exists in the system.
var ErrUserIsExist = errors.New("user is exist")

// ErrInvalidData indicates that the provided data is invalid or malformed.
var ErrInvalidData = errors.New("invalid data")
