package handlers

import (
	"net/http"
	"testing"
)

func TestClientIP(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.RemoteAddr = "127.0.0.1:4567"
	if got := clientIP(req); got != "127.0.0.1" {
		t.Fatalf("clientIP without proxy = %q", got)
	}

	req.Header.Set("X-Forwarded-For", "203.0.113.9, 10.0.0.1")
	if got := clientIP(req); got != "203.0.113.9" {
		t.Fatalf("clientIP with forwarded header = %q", got)
	}
}

func TestBrowserName(t *testing.T) {
	if got := browserName("Mozilla/5.0 Chrome/126.0 Safari/537.36"); got != "Chrome" {
		t.Fatalf("expected Chrome, got %q", got)
	}
	if got := browserName("Mozilla/5.0 Firefox/126.0"); got != "Firefox" {
		t.Fatalf("expected Firefox, got %q", got)
	}
}

func TestShouldSkipActivityAudit(t *testing.T) {
	if !shouldSkipActivityAudit("/api/me") {
		t.Fatal("/api/me should be skipped to avoid noisy background identity checks")
	}
	if shouldSkipActivityAudit("/api/submissions") {
		t.Fatal("workflow endpoints should be activity-audited")
	}
}
