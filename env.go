package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func readEnv() *Env {
	cfg := &Env{}
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	val := os.Getenv("NOTION_SECRET")
	if err != nil {
		log.Fatal("Error reading NOTION_SECRET from .env file")
	}
	cfg.NotionSecret = val
	return cfg
}
