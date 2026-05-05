package password

import (
	"testing"
)

func TestHashAndVerify(t *testing.T) {
	pwd := "Test1234!"
	hash, err := Hash(pwd)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}
	if hash == "" {
		t.Fatal("Hash returned empty string")
	}
	if !Verify(pwd, hash) {
		t.Fatal("Verify failed for correct password")
	}
	if Verify("WrongPassword1!", hash) {
		t.Fatal("Verify should fail for wrong password")
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		pwd  string
		want bool
	}{
		{"Abc1!", false},       // too short (5)
		{"Abcdef1!", true},     // 8 chars, ok
		{"abcdefgh", true},     // 8 chars, ok (no strength check)
		{"12345678", true},     // 8 chars, ok
	}
	for _, tt := range tests {
		err := ValidatePassword(tt.pwd)
		if tt.want && err != nil {
			t.Errorf("ValidatePassword(%q) = %v, want nil", tt.pwd, err)
		}
		if !tt.want && err == nil {
			t.Errorf("ValidatePassword(%q) = nil, want error", tt.pwd)
		}
	}
}

func TestValidatePasswordStrong(t *testing.T) {
	tests := []struct {
		pwd  string
		want bool
	}{
		{"short", false},       // too short
		{"abcdefgh", false},    // no uppercase, digit, special
		{"ABCDEFGH", false},    // no lowercase, digit, special
		{"Test1234", false},    // no special
		{"test1234!", false},   // no uppercase
		{"TEST1234!", false},   // no lowercase
		{"Test!!!!", false},    // no digit
		{"Test1234!", true},    // valid
		{"Abc1!xyz", true},     // valid (lowercase comes later after special)
	}
	for _, tt := range tests {
		err := ValidatePasswordStrong(tt.pwd)
		if tt.want && err != nil {
			t.Errorf("ValidatePasswordStrong(%q) = %v, want nil", tt.pwd, err)
		}
		if !tt.want && err == nil {
			t.Errorf("ValidatePasswordStrong(%q) = nil, want error", tt.pwd)
		}
	}
}
