package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MatyD356/vimGame/internals/cache"
	"github.com/MatyD356/vimGame/internals/env"
	"github.com/MatyD356/vimGame/internals/handlers"
)

type Config struct {
	Port      string
	Env       *env.Env
	NotionApi Client
	Cache     *cache.Cache
}

func main() {
	envCfg, err := env.ReadEnv()
	if err != nil {
		fmt.Println("Error reading environment variables:", err)
		return
	}
	cfg := &Config{
		Env:       envCfg,
		NotionApi: NewClient(10*time.Second, 10*time.Second),
		Cache:     cache.NewCache(),
	}

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/health", handlers.HandleHealt)
	serverMux.Handle("/notion/{databaseId}", http.HandlerFunc(cfg.GetDatabase))

	corsWrappedMux := corsMiddleware(serverMux)
	httpServer := http.Server{
		Addr:              ":" + cfg.Env.Port,
		Handler:           corsWrappedMux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	fmt.Println("Starting server on port", httpServer.Addr)
	err = httpServer.ListenAndServe()

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

}
