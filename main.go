package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MatyD356/vimGame/internals/cache"
	"github.com/MatyD356/vimGame/internals/config"
	"github.com/MatyD356/vimGame/internals/env"
	"github.com/MatyD356/vimGame/internals/handlers"
	"github.com/MatyD356/vimGame/internals/middleware"
)

func main() {
	envCfg, err := env.ReadEnv()
	if err != nil {
		fmt.Println("Error reading environment variables:", err)
		return
	}
	cfg := &config.Config{
		Env:   envCfg,
		Cache: cache.NewCache(),
		HttpClient: &http.Client{
			Timeout: time.Duration(5) * time.Second,
		},
	}

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/health", handlers.HandleHealt)
	serverMux.HandleFunc("/training", handlers.HandleGetTraining)

	corsWrappedMux := middleware.Cors(serverMux)
	depedencyWrappedMux := middleware.DependencyInjection(corsWrappedMux, cfg)
	httpServer := http.Server{
		Addr:              ":" + cfg.Env.Port,
		Handler:           depedencyWrappedMux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	fmt.Println("Starting server on port", httpServer.Addr)
	err = httpServer.ListenAndServe()

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

}
