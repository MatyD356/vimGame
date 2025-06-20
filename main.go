package main

import (
	"fmt"
	"net/http"
	"time"

	env "github.com/MatyD356/vimGame/internals/env"
	"github.com/MatyD356/vimGame/internals/handlers"
)

type Config struct {
	Port   string
	Env    *env.Env
	Client *http.Client
}

func main() {
	// Env
	envCfg, err := env.ReadEnv()
	if err != nil {
		if err == env.ErrMissingNotionSecret {
			fmt.Println("Error: NOTION_SECRET is not set in the environment")
			return
		}
		fmt.Println("Error reading environment variables:", err)
		return
	}
	cfg := &Config{
		Env:    envCfg,
		Client: NewHttpClient(),
		Port:   ":8080",
	}
	fmt.Println("Notion Secret:", cfg.Env.NotionSecret)

	// Server
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/health", handlers.HandleHealt)

	httpServer := http.Server{
		Addr:              ":8080",
		Handler:           serverMux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	err = httpServer.ListenAndServe()

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

}
