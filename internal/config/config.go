package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Version int     `json:"version"`
	Server  Server  `json:"server"`
	Routes  []Route `json:"routes"`
}

type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Route struct {
	Method   string   `json:"method"`
	Path     string   `json:"path"`
	Response Response `json:"response"`
}

type Response struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    json.RawMessage   `json:"body"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()

	var cfg Config

	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("decode config: unexpected data after json")
		}
		return nil, fmt.Errorf("decode config: trailing data: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg.Version != 1 {
		return fmt.Errorf("unsupported config version: %d", cfg.Version)
	}

	if strings.TrimSpace(cfg.Server.Host) == "" {
		return fmt.Errorf("server port must be between 1 and 65535")
	}
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf(("server port must be between 1 and 65535"))
	}

	seenRoutes := make(map[string]bool)

	for i, route := range cfg.Routes {
		if !isSupportedMethod(route.Method) {
			return fmt.Errorf("route %d: unsupported HTTP method %q", i, route.Method)
		}

		if route.Path == "" || route.Path[0] != '/' {
			return fmt.Errorf("route %d: path must start with /", i)
		}

		if strings.ContainsAny(route.Path, "?#") {
			return fmt.Errorf("route %d: path must not contain query or fragment", i)
		}

		key := route.Method + " " + route.Path

		if seenRoutes[key] {
			return fmt.Errorf("route %d: duplicate route %s", i, key)
		}

		seenRoutes[key] = true

		if route.Response.Status < 200 || route.Response.Status > 599 {
			return fmt.Errorf("route %d: Invalid HTTP status", i)
		}
		if len(route.Response.Body) == 0 || !json.Valid(route.Response.Body) {
			return fmt.Errorf("route %d: response body must contain valid JSON", i)
		}

		for name := range route.Response.Headers {
			if !isValidHeaderName(name) {
				return fmt.Errorf("route %d: invalid header name %q", i, name)
			}
		}
	}

	return nil
}

func isSupportedMethod(method string) bool {
	switch method {
	case http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions:
		return true
	default:
		return false
	}
}

func isValidHeaderName(name string) bool {
	if name == "" {
		return false
	}

	for _, char := range name {
		if !isHTTPTokenChar(char) {
			return false
		}
	}
	return true
}

func isHTTPTokenChar(char rune) bool {
	if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' {
		return true
	}

	return strings.ContainsRune("!#$%&'*+-.^_`|~", char)
}
