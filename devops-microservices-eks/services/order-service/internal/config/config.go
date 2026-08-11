package config

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	Port                  string
	DatabaseURL           string
	KafkaBootstrapServers string
	OrderCreatedTopic     string
	Env                   string
}

// Load reads config from env vars and fails fast (exits the process) if
// something required is missing in production — same intent as the Auth
// (Node.js) service's env validation and the Notification (Spring Boot)
// service's EnvValidationConfig. Crash loudly at startup, not confusingly
// on the first request that touches the missing value.
func Load() *Config {
	cfg := &Config{
		Port:                  getEnv("PORT", "8082"),
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		OrderCreatedTopic:     getEnv("ORDER_CREATED_TOPIC", "order.created"),
		Env:                   getEnv("APP_ENV", "development"),
	}

	if cfg.Env == "production" {
		missing := []string{}
		if cfg.DatabaseURL == "" {
			missing = append(missing, "DATABASE_URL")
		}
		if len(missing) > 0 {
			slog.Error("env_validation_failed", "missing", missing)
			fmt.Fprintf(os.Stderr, "missing required environment variables: %v\n", missing)
			os.Exit(1)
		}
		slog.Info("env_validation_passed")
	} else if cfg.DatabaseURL == "" {
		// Local dev fallback so `go run .` works without a full .env setup.
		// This path is intentionally never reachable in production (guarded above).
		cfg.DatabaseURL = "postgres://postgres:postgres@localhost:5432/orders?sslmode=disable"
		slog.Warn("using_default_database_url_dev_only")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
