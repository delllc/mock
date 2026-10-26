package config

import (
	"encoding/json"
	"strings"
	"testing"
)

func validConfig() *Config {
	return &Config{
		Version: 1,
		Server: Server{
			Host: "127.0.0.1",
			Port: 3001,
		},
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/api/products",
				Response: Response{
					Status: 200,
					Body:   json.RawMessage(`[]`),
				},
			},
		},
	}
}

func TestValidateAcceptsValidConfig(t *testing.T) {
	cfg := validConfig()

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
}

func TestValidateRejectsDuplicateRoutes(t *testing.T) {
	cfg := validConfig()
	cfg.Routes = append(cfg.Routes, cfg.Routes[0])

	err := cfg.Validate()

	if err == nil || !strings.Contains(err.Error(), "duplicate route") {
		t.Fatalf("expected duplicate route error, got: %v", err)
	}
}

func TestValidateRejectsInvalidMethod(t *testing.T) {
	cfg := validConfig()
	cfg.Routes[0].Method = "INVALID"

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid method error")
	}
}

func TestValidateRejectsInvalidHeader(t *testing.T) {
	cfg := validConfig()
	cfg.Routes[0].Response.Headers = map[string]string{
		"Bad Header": "value",
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid header error")
	}
}

func TestValidateAllowsSamePathWithDifferentMethods(t *testing.T) {
	cfg := validConfig()

	secondRoute := cfg.Routes[0]
	secondRoute.Method = "POST"

	cfg.Routes = append(cfg.Routes, secondRoute)

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
}
