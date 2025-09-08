package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "mysql-fintech-app/internal/middleware"
    "mysql-fintech-app/internal/services"
)

type TransactionHandler struct { tx *services.TransactionService }
func NewTransactionHandler(s *services.TransactionService) *TransactionHandler { return &TransactionHandler{tx: s} }

type amountReq struct { Amount float64 `json:"amount"` }
type transferReq struct { ToUserID int64 `json:"to_user_id"`; Amount float64 `json:"amount"` }

func (h *TransactionHandler) Credit(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(int64)
    var req amountReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "bad_request", http.StatusBadRequest); return }
    if err := h.tx.Credit(r.Context(), userID, req.Amount); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    w.WriteHeader(http.StatusNoContent)
}

func (h *TransactionHandler) Debit(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(int64)
    var req amountReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "bad_request", http.StatusBadRequest); return }
    if err := h.tx.Debit(r.Context(), userID, req.Amount); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    w.WriteHeader(http.StatusNoContent)
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
    fromID := r.Context().Value(middleware.UserIDKey).(int64)
    var req transferReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "bad_request", http.StatusBadRequest); return }
    if err := h.tx.Transfer(r.Context(), fromID, req.ToUserID, req.Amount); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    w.WriteHeader(http.StatusNoContent)
}

func (h *TransactionHandler) History(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(int64)
    q := r.URL.Query()
    limit, _ := strconv.Atoi(q.Get("limit")); if limit <= 0 { limit = 20 }
    offset, _ := strconv.Atoi(q.Get("offset")); if offset < 0 { offset = 0 }
    items, err := h.tx.History(r.Context(), userID, limit, offset)
    if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
    writeJSON(w, http.StatusOK, items)
}
