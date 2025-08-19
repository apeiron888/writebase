package infrastructure

import "testing"

func TestPasswordHashAndVerify(t *testing.T) {
	svc := NewPasswordService()
	hash, err := svc.HashPassword("MyP@ssw0rd!")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if !svc.VerifyPassword(hash, "MyP@ssw0rd!") {
		t.Fatalf("expected verify ok")
	}
	if svc.VerifyPassword(hash, "bad") {
		t.Fatalf("expected verify fail")
	}
}

func TestIsPasswordStrong(t *testing.T) {
	svc := NewPasswordService()
	if svc.IsPasswordStrong("weak") {
		t.Fatalf("expected weak false")
	}
	if !svc.IsPasswordStrong("Str0ng!Pwd") {
		t.Fatalf("expected strong true")
	}
}
