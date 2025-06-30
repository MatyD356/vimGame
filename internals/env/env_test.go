package env

import (
	"os"
	"testing"
)

func TestReadEnv(t *testing.T) {
	// Create a temporary .env file
	envContent := "NOTION_SECRET=test123\nNOTION_DB_ID=123\nPORT=8080\n"
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write .env file: %v", err)
	}
	defer os.Remove(".env")
	defer os.Unsetenv("NOTION_SECRET") // Clean up environment variable after test

	cfg, err := ReadEnv()
	if err != nil {
		t.Fatalf("rRadEnv() returned an error: %v", err)
	}

	if cfg.NotionSecret != "test123" {
		t.Errorf("Expected NotionSecret to be 'test123', got '%s'", cfg.NotionSecret)
	}
}

func TestReadEnv_MissingSecret(t *testing.T) {
	// Write .env file without NOTION_SECRET
	err := os.WriteFile(".env", []byte("OTHER_VAR=abc\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write .env file: %v", err)
	}
	defer os.Remove(".env")

	_, err = ReadEnv()
	expectedErr := "NOTION_SECRET is not set in the environment"
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected %s, got %v", expectedErr, err)
	}
}
