package domain

import "testing"

func TestSetPassword_And_CheckPassword(t *testing.T) {
	u := &User{}
	if err := u.SetPassword("password123"); err != nil {
		t.Fatalf("SetPassword failed: %v", err)
	}
	if u.PasswordHash == "" {
		t.Fatal("PasswordHash should not be empty")
	}
	if u.PasswordHash == "password123" {
		t.Fatal("PasswordHash should not be plaintext")
	}
	if !u.CheckPassword("password123") {
		t.Error("CheckPassword should return true for correct password")
	}
	if u.CheckPassword("wrongpassword") {
		t.Error("CheckPassword should return false for wrong password")
	}
}
