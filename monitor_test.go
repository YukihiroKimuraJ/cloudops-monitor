package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckOnce_Status200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	client := &http.Client{Timeout: 200 * time.Millisecond}

	code, _, err := CheckOnce(context.Background(), client, srv.URL)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, code)
	}
}

func TestCheckOnce_Status500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	client := &http.Client{Timeout: 200 * time.Millisecond}

	code, _, err := CheckOnce(context.Background(), client, srv.URL)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, code)
	}
}

func TestCheckOnce_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond) // わざと遅延
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	// タイムアウトを短く設定して確実にエラーを起こす
	client := &http.Client{Timeout: 30 * time.Millisecond}

	_, _, err := CheckOnce(context.Background(), client, srv.URL)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
	// タイムアウトエラーかどうか確認
	// net/http のタイムアウトは os.IsTimeout または文字列チェックで判定
	if !isTimeoutError(err) {
		t.Fatalf("expected timeout error, got: %v", err)
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"example.com", "https://example.com", false},
		{"http://example.com", "http://example.com", false},
		{"https://example.com", "https://example.com", false},
		{"", "", true},
		{"  ", "", true},
	}

	for _, tt := range tests {
		got, err := normalizeURL(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("normalizeURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if got != tt.expected {
			t.Errorf("normalizeURL(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// isTimeoutError checks if the error is a timeout error
func isTimeoutError(err error) bool {
	type timeoutError interface {
		Timeout() bool
	}
	te, ok := err.(timeoutError)
	return ok && te.Timeout()
}
