package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL         string
	OpenAIAPIKey        string
	DigiKeyClientID     string
	DigiKeyClientSecret string
	MouserSearchAPIKey  string
	PollIntervalMs      int
}

func Load() *Config {
	// Try loading .env file, ignore if it doesn't exist (e.g., in production)
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pollIntervalMs := 5000 // default
	if val := os.Getenv("POLL_INTERVAL_MS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			pollIntervalMs = parsed
		}
	}

	return &Config{
		DatabaseURL:         dbURL,
		OpenAIAPIKey:        os.Getenv("OPENAI_API_KEY"),
		DigiKeyClientID:     os.Getenv("DIGIKEY_CLIENT_ID"),
		DigiKeyClientSecret: os.Getenv("DIGIKEY_CLIENT_SECRET"),
		MouserSearchAPIKey:  os.Getenv("MOUSER_SEARCH_API_KEY"),
		PollIntervalMs:      pollIntervalMs,
	}
}
