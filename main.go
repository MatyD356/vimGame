package main

import (
	"fmt"
	"net/http"

	"github.com/MatyD356/vimGame/internals/handlers"
)

type Config struct {
	Port   string
	Env    *Env
	Client *http.Client
}

type Env struct {
	NotionSecret string
}

func main() {
	// Env
	cfg := &Config{
		Env:    readEnv(),
		Client: NewHttpClient(),
		Port:   ":8080",
	}
	fmt.Println("Notion Secret:", cfg.Env.NotionSecret)

	// Server
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/health", handlers.HandleHealt)

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	httpServer.ListenAndServe()

}
