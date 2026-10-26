package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

	if cfg.Server.Host == "" {
		return fmt.Errorf("server port must be between 1 and 65535")
	}
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf(("server port must be between 1 and 65535"))
	}

	for i, route := range cfg.Routes {
		if route.Method == "" {
			return fmt.Errorf("route %d: method is required", i)
		}

		if route.Path == "" || route.Path[0] != '/' {
			return fmt.Errorf("route %d: path must start with /", i)
		}

		if route.Response.Status < 100 || route.Response.Status > 599 {
			return fmt.Errorf("route %d: Invalid HTTP status", i)
		}
		if len(route.Response.Body) == 0 || !json.Valid(route.Response.Body) {
			return fmt.Errorf("route %d: response body must contain valid JSON", i)
		}
	}

	return nil
}
