package repositories

import (
    "context"
    "database/sql"
    "mysql-fintech-app/internal/models"
)

type UserRepo struct { db *sql.DB }

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, u *models.User) (int64, error) {
    q := "INSERT INTO users (username, email, password_hash, role) VALUES (?,?,?,?)"
    res, err := r.db.ExecContext(ctx, q, u.Username, u.Email, u.PasswordHash, u.Role)
    if err != nil { return 0, err }
    id, _ := res.LastInsertId()
    return id, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    q := "SELECT id, username, email, password_hash, role FROM users WHERE email=?"
    row := r.db.QueryRowContext(ctx, q, email)
    u := models.User{}
    if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role); err != nil {
        if err == sql.ErrNoRows { return nil, nil }
        return nil, err
    }
    return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
    q := "SELECT id, username, email, password_hash, role FROM users WHERE id=?"
    row := r.db.QueryRowContext(ctx, q, id)
    u := models.User{}
    if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role); err != nil {
        if err == sql.ErrNoRows { return nil, nil }
        return nil, err
    }
    return &u, nil
}
