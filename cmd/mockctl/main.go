package main

import (
	"fmt"
	"os"

	"mock-platform/internal/config"
)

func main() {
	cfg, err := config.Load("examples/shop.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Routes: %d\n", len(cfg.Routes))

	for _, route := range cfg.Routes {
		fmt.Printf("%s %s\n", route.Method, route.Path)
	}
}
