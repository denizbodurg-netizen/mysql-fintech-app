package api

import (
    "database/sql"
    "net/http"

    "github.com/rs/zerolog"
    "mysql-fintech-app/config"
    "mysql-fintech-app/internal/api/handlers"
    "mysql-fintech-app/internal/middleware"
    "mysql-fintech-app/internal/services"
)

func chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
    for i := len(mws)-1; i >= 0; i-- { h = mws[i](h) }
    return h
}

func NewRouter(log zerolog.Logger, db *sql.DB, cfg *config.Config) http.Handler {
    userSvc := services.NewUserService(cfg, db)
    txSvc   := services.NewTransactionService(db)

    authH := handlers.NewAuthHandler(userSvc)
    balH  := handlers.NewBalanceHandler(txSvc)
    txH   := handlers.NewTransactionHandler(txSvc)

    mux := http.NewServeMux()
    // Public
    mux.Handle("/api/v1/auth/register", chain(http.HandlerFunc(authH.Register)))
    mux.Handle("/api/v1/auth/login",    chain(http.HandlerFunc(authH.Login)))

    // Protected
    authMW := middleware.AuthJWT(cfg.JWTSecret)
    mux.Handle("/api/v1/balances/current",      chain(http.HandlerFunc(balH.Current), authMW))
    mux.Handle("/api/v1/transactions/credit",   chain(http.HandlerFunc(txH.Credit),   authMW))
    mux.Handle("/api/v1/transactions/debit",    chain(http.HandlerFunc(txH.Debit),    authMW))
    mux.Handle("/api/v1/transactions/transfer", chain(http.HandlerFunc(txH.Transfer), authMW))
    mux.Handle("/api/v1/transactions/history",  chain(http.HandlerFunc(txH.History),  authMW))

    // Health
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
    return middleware.Logging(log)(mux)
}
