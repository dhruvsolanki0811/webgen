package service

import (
	"testing"
)

func TestTokenService_GenerateAndValidate(t *testing.T) {
	svc := NewTokenService("test-secret-at-least-32-chars-long!")

	token, err := svc.GenerateAccess("user123")
	if err != nil {
		t.Fatalf("GenerateAccess failed: %v", err)
	}
	if token == "" {
		t.Fatal("token should not be empty")
	}

	userID, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if userID != "user123" {
		t.Errorf("expected user123, got %s", userID)
	}
}

func TestTokenService_InvalidToken(t *testing.T) {
	svc := NewTokenService("test-secret-at-least-32-chars-long!")

	_, err := svc.ValidateToken("garbage-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestTokenService_WrongSecret(t *testing.T) {
	svc1 := NewTokenService("secret-one-at-least-32-chars-long!")
	svc2 := NewTokenService("secret-two-at-least-32-chars-long!")

	token, _ := svc1.GenerateAccess("user123")
	_, err := svc2.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for token signed with different secret")
	}
}

func TestTokenService_RefreshToken(t *testing.T) {
	svc := NewTokenService("test-secret-at-least-32-chars-long!")

	token, err := svc.GenerateRefresh("user456")
	if err != nil {
		t.Fatalf("GenerateRefresh failed: %v", err)
	}

	userID, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if userID != "user456" {
		t.Errorf("expected user456, got %s", userID)
	}
}
