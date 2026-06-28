package ws

import (
	"net/http"
	"testing"
)

func requestWithOrigin(origin string) *http.Request {
	req := &http.Request{Header: http.Header{}}
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	return req
}

func TestCheckOriginAllowsDefaultLocalDevelopmentOrigins(t *testing.T) {
	t.Setenv(allowedOriginsEnv, "")

	allowed := []string{
		"http://localhost:3000",
		"http://localhost:5173",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:5173",
	}

	for _, origin := range allowed {
		if !checkOrigin(requestWithOrigin(origin)) {
			t.Fatalf("expected default origin %q to be allowed", origin)
		}
	}
}

func TestCheckOriginAllowsConfiguredOrigin(t *testing.T) {
	t.Setenv(allowedOriginsEnv, "https://app.example.com, https://admin.example.com")

	if !checkOrigin(requestWithOrigin("https://admin.example.com")) {
		t.Fatal("expected configured origin to be allowed")
	}
}

func TestCheckOriginRejectsUnconfiguredOrigin(t *testing.T) {
	t.Setenv(allowedOriginsEnv, "https://app.example.com")

	if checkOrigin(requestWithOrigin("https://evil.example.com")) {
		t.Fatal("expected unconfigured origin to be rejected")
	}
}

func TestCheckOriginRejectsMissingOrigin(t *testing.T) {
	t.Setenv(allowedOriginsEnv, "https://app.example.com")

	if checkOrigin(requestWithOrigin("")) {
		t.Fatal("expected missing origin to be rejected")
	}
}

func TestCheckOriginRejectsMalformedOrigin(t *testing.T) {
	t.Setenv(allowedOriginsEnv, "https://app.example.com")

	if checkOrigin(requestWithOrigin("://not-a-valid-origin")) {
		t.Fatal("expected malformed origin to be rejected")
	}
}
