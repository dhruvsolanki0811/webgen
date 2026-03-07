package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrNotFound(t *testing.T) {
	err := ErrNotFound("user")
	if err.Code != 404 {
		t.Errorf("expected 404, got %d", err.Code)
	}
	if err.Message != "user not found" {
		t.Errorf("expected 'user not found', got %s", err.Message)
	}
}

func TestErrInternal_WrapsError(t *testing.T) {
	original := fmt.Errorf("db connection failed")
	err := ErrInternal(original)
	if err.Code != 500 {
		t.Errorf("expected 500, got %d", err.Code)
	}
	if err.Message != "internal error" {
		t.Errorf("expected 'internal error', got %s", err.Message)
	}
	if !errors.Is(err, original) {
		t.Error("expected Unwrap to return original error")
	}
}

func TestAppError_ErrorString(t *testing.T) {
	err := ErrInternal(fmt.Errorf("something broke"))
	if err.Error() != "internal error: something broke" {
		t.Errorf("unexpected error string: %s", err.Error())
	}
}

func TestAppError_ErrorsAs(t *testing.T) {
	err := ErrBadRequest("bad input")
	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Error("errors.As should work with AppError")
	}
}
