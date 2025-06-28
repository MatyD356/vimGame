package env

import (
	"errors"
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
		return nil, errors.New("NOTION_SECRET is not set in the environment")
	}
	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("PORT is not set in the environment")
	}

	notionDbId := os.Getenv("NOTION_DB_ID")
	if notionDbId == "" {
		return nil, errors.New("NOTION_DB_ID is not set in the environment")
	}

	return &Env{
		NotionSecret: secret,
		NotionDbId:   notionDbId,
		Port:         port,
	}, nil
}
