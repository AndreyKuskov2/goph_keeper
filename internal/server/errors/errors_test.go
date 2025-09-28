package postgres

import (
	"errors"
	"testing"
)

func TestCustomErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrUserIsExist",
			err:      ErrUserIsExist,
			expected: "user is exist",
		},
		{
			name:     "ErrInvalidData",
			err:      ErrInvalidData,
			expected: "invalid data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Expected error message '%s', got '%s'", tt.expected, tt.err.Error())
			}
		})
	}
}

func TestErrorComparison(t *testing.T) {
	// Test that errors.Is works correctly
	if !errors.Is(ErrUserIsExist, ErrUserIsExist) {
		t.Error("errors.Is should return true for same error")
	}

	if errors.Is(ErrUserIsExist, ErrInvalidData) {
		t.Error("errors.Is should return false for different errors")
	}

	// Test wrapped errors
	wrappedErr := errors.New("wrapped: " + ErrUserIsExist.Error())
	if errors.Is(wrappedErr, ErrUserIsExist) {
		t.Error("errors.Is should return false for wrapped error without proper wrapping")
	}
}

func TestErrorUniqueness(t *testing.T) {
	// Test that errors are unique instances
	if ErrUserIsExist == ErrInvalidData {
		t.Error("Different errors should not be equal")
	}

	// Test that creating new instances with same message are different
	newErr := errors.New("user is exist")
	if ErrUserIsExist == newErr {
		t.Error("Different error instances should not be equal even with same message")
	}
}
