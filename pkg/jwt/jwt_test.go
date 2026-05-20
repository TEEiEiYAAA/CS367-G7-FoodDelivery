package jwt

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	t.Run("success - returns non-empty token", func(t *testing.T) {
		token, err := GenerateToken("testuser", "customer")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected non-empty token, got empty string")
		}
	})
}

func TestValidateToken(t *testing.T) {
	t.Run("success - valid token returns correct claims", func(t *testing.T) {
		token, _ := GenerateToken("testuser", "customer")

		claims, err := ValidateToken(token)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if claims["username"] != "testuser" {
			t.Errorf("expected username 'testuser', got %v", claims["username"])
		}
		if claims["role"] != "customer" {
			t.Errorf("expected role 'customer', got %v", claims["role"])
		}
	})

	t.Run("failure - invalid token returns error", func(t *testing.T) {
		_, err := ValidateToken("invalid.token.string")
		if err == nil {
			t.Error("expected error for invalid token, got nil")
		}
	})

	t.Run("failure - empty token returns error", func(t *testing.T) {
		_, err := ValidateToken("")
		if err == nil {
			t.Error("expected error for empty token, got nil")
		}
	})
}
