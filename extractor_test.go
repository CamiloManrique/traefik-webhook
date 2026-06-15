package urlextractor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	urlextractor "github.com/CamiloManrique/traefik-url-extractor"
)

func TestURLExtractor(t *testing.T) {
	cfg := urlextractor.CreateConfig()
	cfg.Regex = `/flows/(?P<flowId>[a-f0-9-]+)/components/(?P<componentId>[a-zA-Z0-9_-]+)`
	cfg.Headers["X-Flow-Id"] = "flowId"
	cfg.Headers["X-Component-Id"] = "componentId"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := urlextractor.New(ctx, next, cfg, "test-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/flows/abc123/components/my-trigger", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-Flow-Id", "abc123")
	assertHeader(t, req, "X-Component-Id", "my-trigger")
}

func TestURLExtractorNoMatch(t *testing.T) {
	cfg := urlextractor.CreateConfig()
	cfg.Regex = `/flows/(?P<flowId>[a-f0-9-]+)`
	cfg.Headers["X-Flow-Id"] = "flowId"

	ctx := context.Background()

	var nextCalled bool
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		nextCalled = true
	})

	handler, err := urlextractor.New(ctx, next, cfg, "test-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/other/path", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	if !nextCalled {
		t.Error("next handler should have been called on no match")
	}
	if req.Header.Get("X-Flow-Id") != "" {
		t.Error("header should not be set on no match")
	}
}

func TestURLExtractorInvalidConfig(t *testing.T) {
	t.Run("empty regex", func(t *testing.T) {
		cfg := urlextractor.CreateConfig()
		cfg.Headers["X-Id"] = "id"
		_, err := urlextractor.New(context.Background(), http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}), cfg, "test")
		if err == nil {
			t.Error("expected error for empty regex")
		}
	})

	t.Run("empty headers", func(t *testing.T) {
		cfg := urlextractor.CreateConfig()
		cfg.Regex = `/(?P<id>[a-z]+)`
		_, err := urlextractor.New(context.Background(), http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}), cfg, "test")
		if err == nil {
			t.Error("expected error for empty headers")
		}
	})

	t.Run("invalid regex", func(t *testing.T) {
		cfg := urlextractor.CreateConfig()
		cfg.Regex = `(?P<id>[`
		cfg.Headers["X-Id"] = "id"
		_, err := urlextractor.New(context.Background(), http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}), cfg, "test")
		if err == nil {
			t.Error("expected error for invalid regex")
		}
	})
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()
	if req.Header.Get(key) != expected {
		t.Errorf("header %s: got %q, want %q", key, req.Header.Get(key), expected)
	}
}
