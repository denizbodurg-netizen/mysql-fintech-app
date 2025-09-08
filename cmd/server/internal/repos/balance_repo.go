package repositories

import (
    "context"
    "database/sql"
    "errors"
)

type BalanceRepo struct { db *sql.DB }

func NewBalanceRepo(db *sql.DB) *BalanceRepo { return &BalanceRepo{db: db} }

// delta kadar bakiyeyi artır/azalt (transaction içinde güvenli)
func (r *BalanceRepo) AdjustBalance(ctx context.Context, tx *sql.Tx, userID int64, delta float64) (float64, error) {
    // Eğer satır yoksa oluştur (PRIMARY KEY user_id; varsa yok sayar)
    if _, err := tx.ExecContext(ctx,
        "INSERT IGNORE INTO balances (user_id, amount) VALUES (?, 0)",
        userID,
    ); err != nil {
        return 0, err
    }
    // Satırı kilitle, mevcut bakiyeyi oku
    var current float64
    if err := tx.QueryRowContext(ctx,
        "SELECT amount FROM balances WHERE user_id=? FOR UPDATE",
        userID,
    ).Scan(&current); err != nil {
        return 0, err
    }

    newAmount := current + delta
    if newAmount < 0 { return current, errors.New("insufficient_funds") }

    if _, err := tx.ExecContext(ctx,
        "UPDATE balances SET amount=? WHERE user_id=?",
        newAmount, userID,
    ); err != nil {
        return current, err
    }
    return newAmount, nil
}

func (r *BalanceRepo) GetCurrent(ctx context.Context, userID int64) (float64, error) {
    var amt float64
    err := r.db.QueryRowContext(ctx, "SELECT amount FROM balances WHERE user_id=?", userID).Scan(&amt)
    if err == sql.ErrNoRows { return 0, nil }
    return amt, err
}
