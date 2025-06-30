package middleware

import (
	"context"
	"net/http"

	"github.com/MatyD356/vimGame/internals/config"
)

type contextKey string

const ConfigKey contextKey = "config"

func DependencyInjection(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ConfigKey, cfg)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
