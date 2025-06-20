package env

import (
	"os"

	"github.com/joho/godotenv"
)

func ReadEnv() (*Env, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	secret := os.Getenv("NOTION_SECRET")
	if secret == "" {
		return nil, ErrMissingNotionSecret
	}

	return &Env{NotionSecret: secret}, nil
}
