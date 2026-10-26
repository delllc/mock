package mockengine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mock-platform/internal/config"
)

func TestEngineReturnsConfiguredResponse(t *testing.T) {
	cfg := &config.Config{
		Routes: []config.Route{
			{
				Method: "GET",
				Path:   "/api/products",
				Response: config.Response{
					Status: 200,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: json.RawMessage(
						`[{"id":1,"name":"MacBook"}]`,
					),
				},
			},
		},
	}

	engine := New(cfg)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/products",
		nil,
	)
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", got)
	}

	expectedBody := `[{"id":1,"name":"MacBook"}]`
	if rec.Body.String() != expectedBody {
		t.Fatalf(
			"expected body %s, got %s",
			expectedBody,
			rec.Body.String(),
		)
	}
}

func TestEngineReturnsNotFoundForUnknownRoute(t *testing.T) {
	engine := New(&config.Config{})

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/unknown",
		nil,
	)
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}

func TestEngineMatchesHTTPMethod(t *testing.T) {
	cfg := &config.Config{
		Routes: []config.Route{
			{
				Method: "GET",
				Path:   "/api/products",
				Response: config.Response{
					Status: http.StatusOK,
					Body:   json.RawMessage(`[]`),
				},
			},
		},
	}

	engine := New(cfg)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/products",
		nil,
	)
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}
