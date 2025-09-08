package models

type Balance struct {
    UserID        int64   `json:"user_id"`
    Amount        float64 `json:"amount"`
    LastUpdatedAt string  `json:"last_updated_at"`
}
