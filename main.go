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

var whitelist = []string{
	"http://localhost:5137",
}

func isAllowedOrigin(origin string) bool {
	for _, allowed := range whitelist {
		if origin == allowed {
			return true
		}
	}
	return false
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
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
		Env:       envCfg,
		NotionApi: NewClient(10*time.Second, 10*time.Second),
		Port:      ":8080",
		Cache:     cache.NewCache(),
	}
	fmt.Println("Notion Secret:", cfg.Env.NotionSecret)

	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/health", handlers.HandleHealt)
	serverMux.Handle("/notion/{databaseId}", corsMiddleware(http.HandlerFunc(cfg.GetDatabase)))

	corsWrappedMux := corsMiddleware(serverMux)
	httpServer := http.Server{
		Addr:              ":8080",
		Handler:           corsWrappedMux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	err = httpServer.ListenAndServe()

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

}
