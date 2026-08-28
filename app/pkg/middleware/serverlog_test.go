package middleware

import (
	"net/http"
	"testing"
)

func TestFilterHeadersHidesAVMWeChatSecret(t *testing.T) {
	headers := http.Header{
		"Content-Type":        {"application/json"},
		"X-Avm-Wechat-Secret": {"bridge-secret"},
	}

	filtered := filterHeaders(headers, map[string]bool{}, hideShowHeaders)

	if got := filtered.Get("X-AVM-WeChat-Secret"); got != "" {
		t.Fatalf("X-AVM-WeChat-Secret leaked into logs: %q", got)
	}
	if got := filtered.Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	if got := headers.Get("X-AVM-WeChat-Secret"); got != "bridge-secret" {
		t.Fatalf("request headers were mutated: %q", got)
	}
}
