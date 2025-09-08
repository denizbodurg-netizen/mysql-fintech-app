package config

import(
	"os"
	"time"
	"log"
)

type Config struct{
	Env string
	DBUrl string
	JWTSecret string
	JWTExpiry time.Duration
}

func Load() *Config{
	cfg := &Config{
		Env:	getEnv("APP_ENV", "development"),
		DBUrl:  getEnv("DATABASE_URL", "appuser:apppass@tcp(localhost:3306)/fintech?parseTime=true&loc=Local&multiStatements=true"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret"),
		JWTExpiry: getDuration("JWT_Expiry", 24*time.Hour),
	}
	if cfg.Env != "development" && cfg.JWTSecret == "dev-secret-change-me" {
		log.Fatal("JWT_SECRET must be set in non-development envs")
	}
	if cfg.DBUrl == "" {
		log.Fatal("DATABASE_URL is required")
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v:= os.Getenv(key); v != "" {return v}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration{
	if v:= os.Getenv(key); v != ""{
		if d, err := time.ParseDuration(v); err == nil {return d}
	}
	return fallback
}