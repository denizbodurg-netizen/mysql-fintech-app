package handlers

import (
    "encoding/json"
    "net/http"
    "mysql-fintech-app/internal/services"
)

type AuthHandler struct { user *services.UserService }
func NewAuthHandler(u *services.UserService) *AuthHandler { return &AuthHandler{user: u} }

type registerReq struct { Username, Email, Password string }
type registerResp struct { ID int64 `json:"id"` }
type loginReq struct { Email, Password string }
type loginResp struct { Token string `json:"token"`; UserID int64 `json:"user_id"` }

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req registerReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "bad_request", http.StatusBadRequest); return }
    id, err := h.user.Register(r.Context(), req.Username, req.Email, req.Password)
    if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
    writeJSON(w, http.StatusCreated, registerResp{ID: id})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req loginReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, "bad_request", http.StatusBadRequest); return }
    token, id, err := h.user.Login(r.Context(), req.Email, req.Password)
    if err != nil { http.Error(w, err.Error(), http.StatusUnauthorized); return }
    writeJSON(w, http.StatusOK, loginResp{Token: token, UserID: id})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _ = json.NewEncoder(w).Encode(v)
}
