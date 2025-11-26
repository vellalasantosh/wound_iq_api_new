package config

import (
    "errors"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    DB_DSN   string
    Port     string
    AppEnv   string
    LogLevel string
}

func Load() (*Config, error) {
    godotenv.Load() // ignore error (optional .env)

    dsn := os.Getenv("DB_DSN")
    port := os.Getenv("PORT")
    env := os.Getenv("APP_ENV")
    logLevel := os.Getenv("LOG_LEVEL")

    if dsn == "" {
        return nil, errors.New("DB_DSN is required")
    }
    if port == "" {
        port = "8080"
    }
    if env == "" {
        env = "development"
    }
    if logLevel == "" {
        logLevel = "info"
    }
    return &Config{
        DB_DSN:   dsn,
        Port:     port,
        AppEnv:   env,
        LogLevel: logLevel,
    }, nil
}
