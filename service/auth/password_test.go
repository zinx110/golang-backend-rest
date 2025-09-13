package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash")
	}
	if hash == "password" {
		t.Error("hash should not be the same as the password")
	}

}

func TestComparePasswords(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}
	if !ComparePasswords(hash, []byte("password")) {
		t.Error("expected passwords to match")
	}
	if ComparePasswords(hash, []byte("wrongpassword")) {
		t.Error("expected passwords to not match")
	}
}
