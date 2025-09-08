package services

import (
    "context"
    "database/sql"
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"

    "mysql-fintech-app/config"
    "mysql-fintech-app/internal/models"
    "mysql-fintech-app/internal/repositories"
)

type UserService struct {
    cfg   *config.Config
    users *repositories.UserRepo
}

func NewUserService(cfg *config.Config, db *sql.DB) *UserService {
    return &UserService{cfg: cfg, users: repositories.NewUserRepo(db)}
}

func (s *UserService) Register(ctx context.Context, username, email, password string) (int64, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil { return 0, err }
    u := &models.User{Username: username, Email: email, PasswordHash: string(hash), Role: "user"}
    return s.users.Create(ctx, u)
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, int64, error) {
    u, err := s.users.GetByEmail(ctx, email)
    if err != nil { return "", 0, err }
    if u == nil { return "", 0, errors.New("invalid_credentials") }
    if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
        return "", 0, errors.New("invalid_credentials")
    }
    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub":  u.ID,
        "role": u.Role,
        "exp":  time.Now().Add(s.cfg.JWTExpiry).Unix(),
    })
    signed, err := tok.SignedString([]byte(s.cfg.JWTSecret))
    return signed, u.ID, err
}
