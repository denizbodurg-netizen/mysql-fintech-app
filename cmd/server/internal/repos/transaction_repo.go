package repositories

import (
    "context"
    "database/sql"
    "mysql-fintech-app/internal/models"
)

type TransactionRepo struct { db *sql.DB }

func NewTransactionRepo(db *sql.DB) *TransactionRepo { return &TransactionRepo{db: db} }

func (r *TransactionRepo) Create(ctx context.Context, tx *sql.Tx, t *models.Transaction) (int64, error) {
    q := "INSERT INTO transactions (from_user_id, to_user_id, amount, type, status) VALUES (?,?,?,?,?)"
    res, err := tx.ExecContext(ctx, q, t.FromUserID, t.ToUserID, t.Amount, t.Type, t.Status)
    if err != nil { return 0, err }
    id, _ := res.LastInsertId()
    return id, nil
}

func (r *TransactionRepo) ListByUser(ctx context.Context, userID int64, limit, offset int) ([]models.Transaction, error) {
    q := `SELECT id, from_user_id, to_user_id, amount, type, status, created_at
          FROM transactions
          WHERE from_user_id=? OR to_user_id=?
          ORDER BY created_at DESC
          LIMIT ? OFFSET ?`
    rows, err := r.db.QueryContext(ctx, q, userID, userID, limit, offset)
    if err != nil { return nil, err }
    defer rows.Close()

    var out []models.Transaction
    for rows.Next() {
        var t models.Transaction
        var from, to sql.NullInt64
        if err := rows.Scan(&t.ID, &from, &to, &t.Amount, &t.Type, &t.Status, &t.CreatedAt); err != nil { return nil, err }
        if from.Valid { v := from.Int64; t.FromUserID = &v }
        if to.Valid   { v := to.Int64;   t.ToUserID   = &v }
        out = append(out, t)
    }
    return out, rows.Err()
}
