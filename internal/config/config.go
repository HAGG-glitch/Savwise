package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	GroqAPIKey  string
	GroqModel   string
	AppEnv      string
}

func Load() Config {
	loadDotEnv(".env")
	return Config{
		AppPort:     env("APP_PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/savwise_ai?sslmode=disable"),
		GroqAPIKey:  env("GROQ_API_KEY", ""),
		GroqModel:   env("GROQ_MODEL", "llama-3.1-8b-instant"),
		AppEnv:      env("APP_ENV", "development"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}
