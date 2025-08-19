package utils

import (
	"regexp"
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	u := &Utils{}
	got := u.GenerateSlug("Hello   World  Go ")
	if got != "Hello-World-Go" {
		t.Fatalf("unexpected slug: %q", got)
	}
}

func TestGenerateShortUUID(t *testing.T) {
	u := &Utils{}
	s := u.GenerateShortUUID()
	if len(s) != 6 {
		t.Fatalf("expected len 6, got %d (%q)", len(s), s)
	}
	if !regexp.MustCompile(`^[0-9a-f]{6}$`).MatchString(s) {
		t.Fatalf("not hex: %q", s)
	}
}
