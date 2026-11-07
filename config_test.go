package sieve

import "testing"

func TestResolveListenAddr(t *testing.T) {
	t.Setenv("PORT", "")
	if got := resolveListenAddr(""); got != ":8080" {
		t.Fatalf("empty: got %q", got)
	}
	if got := resolveListenAddr(":9090"); got != ":9090" {
		t.Fatalf("configured: got %q", got)
	}

	t.Setenv("PORT", "1234")
	if got := resolveListenAddr(":8080"); got != "0.0.0.0:1234" {
		t.Fatalf("PORT overrides yaml: got %q", got)
	}

	t.Setenv("PORT", ":5678")
	if got := resolveListenAddr(":8080"); got != "0.0.0.0:5678" {
		t.Fatalf("PORT with colon: got %q", got)
	}
}
