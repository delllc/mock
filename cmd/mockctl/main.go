package main

import (
	"fmt"
	"net/http"
	"os"

	"mock-platform/internal/config"
	"mock-platform/internal/mockengine"
)

func main() {
	cfg, err := config.Load("examples/shop.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	engine := mockengine.New(cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	fmt.Printf("Mock server listening on http://%s\n", addr)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
