package handlers

import (
    "encoding/json"
    "net/http"
    "mysql-fintech-app/internal/middleware"
    "mysql-fintech-app/internal/services"
)

type BalanceHandler struct { tx *services.TransactionService }
func NewBalanceHandler(s *services.TransactionService) *BalanceHandler { return &BalanceHandler{tx: s} }

func (h *BalanceHandler) Current(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(int64)
    amt, err := h.tx.CurrentBalance(r.Context(), userID)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]any{"amount": amt})
}
