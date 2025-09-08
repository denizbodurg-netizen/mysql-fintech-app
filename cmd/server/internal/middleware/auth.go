package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

type ctxKey string
const UserIDKey ctxKey = "userID"

func AuthJWT(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            auth := r.Header.Get("Authorization")
            parts := strings.SplitN(auth, " ", 2)
            if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
                http.Error(w, "missing_or_invalid_token", http.StatusUnauthorized); return
            }
            tok, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil })
            if err != nil || !tok.Valid { http.Error(w, "invalid_token", http.StatusUnauthorized); return }
            claims, ok := tok.Claims.(jwt.MapClaims)
            if !ok { http.Error(w, "invalid_token", http.StatusUnauthorized); return }
            idf, ok := claims["sub"].(float64)
            if !ok { http.Error(w, "invalid_token", http.StatusUnauthorized); return }
            ctx := context.WithValue(r.Context(), UserIDKey, int64(idf))
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
