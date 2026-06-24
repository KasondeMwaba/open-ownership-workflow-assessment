package services

import (
	"errors"
	"testing"
)

func TestValidateUserIdentity(t *testing.T) {
	name, email, err := validateUserIdentity("  Sam Admin  ", "  SAM@example.COM ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "Sam Admin" || email != "sam@example.com" {
		t.Fatalf("unexpected normalized identity: %q %q", name, email)
	}

	_, _, err = validateUserIdentity("", "not-an-email")
	if !errors.Is(err, ErrInvalidUserInput) {
		t.Fatalf("expected ErrInvalidUserInput, got %v", err)
	}
}
