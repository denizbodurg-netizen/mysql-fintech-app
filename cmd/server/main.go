module mysql-fintech-app

go 1.22

require (
    github.com/go-sql-driver/mysql v1.8.1
    github.com/golang-jwt/jwt/v5 v5.2.1
    github.com/rs/zerolog v1.33.0
    golang.org/x/crypto v0.25.0
)

package main

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "mysql-fintech-app/config"
    "mysql-fintech-app/internal/api"
    "mysql-fintech-app/internal/db"
    "mysql-fintech-app/internal/logger"
)

func main() {
    cfg := config.Load()
    log := logger.New(cfg.Env)

    database, err := db.Connect(cfg.DBUrl)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to connect DB")
    }
    if err := db.Migrate(database); err != nil {
        log.Fatal().Err(err).Msg("failed to run migrations")
    }

    h := api.NewRouter(log, database, cfg)
    srv := &http.Server{Addr: ":8080", Handler: h}

    go func() {
        log.Info().Msg("server starting on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal().Err(err).Msg("server failed")
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Error().Err(err).Msg("server shutdown failed")
    }
    log.Info().Msg("server gracefully stopped")
}

