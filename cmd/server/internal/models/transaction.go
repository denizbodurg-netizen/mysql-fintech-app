package models

type Transaction struct {
    ID         int64   `json:"id"`
    FromUserID *int64  `json:"from_user_id,omitempty"`
    ToUserID   *int64  `json:"to_user_id,omitempty"`
    Amount     float64 `json:"amount"`
    Type       string  `json:"type"`
    Status     string  `json:"status"`
    CreatedAt  string  `json:"created_at"`
}
