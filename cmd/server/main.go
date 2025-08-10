package main

import (
	"complaint-escalator/internal/config"
	"fmt"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	server, err := NewServer(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to start server: %v", err))
	}
	server.Start()
}
