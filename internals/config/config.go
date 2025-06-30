package config

import (
	"net/http"

	"github.com/MatyD356/vimGame/internals/cache"
	"github.com/MatyD356/vimGame/internals/env"
)

type Config struct {
	Env        *env.Env
	Cache      *cache.Cache
	HttpClient *http.Client
}

func Create(env *env.Env, cache *cache.Cache, client *http.Client) *Config {
	return &Config{
		Env:        env,
		Cache:      cache,
		HttpClient: client,
	}
}
