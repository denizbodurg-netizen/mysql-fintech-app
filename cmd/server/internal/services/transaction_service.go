package services

import (
    "context"
    "database/sql"
    "errors"

    "mysql-fintech-app/internal/models"
    "mysql-fintech-app/internal/repositories"
)

type TransactionService struct {
    db    *sql.DB
    bal   *repositories.BalanceRepo
    txRep *repositories.TransactionRepo
}

func NewTransactionService(db *sql.DB) *TransactionService {
    return &TransactionService{db: db, bal: repositories.NewBalanceRepo(db), txRep: repositories.NewTransactionRepo(db)}
}

func (s *TransactionService) Credit(ctx context.Context, userID int64, amount float64) error {
    if amount <= 0 { return errors.New("amount_must_be_positive") }
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()

    if _, err := s.bal.AdjustBalance(ctx, tx, userID, amount); err != nil { return err }
    var from *int64 = nil
    to := &userID
    t := &models.Transaction{FromUserID: from, ToUserID: to, Amount: amount, Type: "credit", Status: "completed"}
    if _, err := s.txRep.Create(ctx, tx, t); err != nil { return err }
    return tx.Commit()
}

func (s *TransactionService) Debit(ctx context.Context, userID int64, amount float64) error {
    if amount <= 0 { return errors.New("amount_must_be_positive") }
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()

    if _, err := s.bal.AdjustBalance(ctx, tx, userID, -amount); err != nil { return err }
    from := &userID
    var to *int64 = nil
    t := &models.Transaction{FromUserID: from, ToUserID: to, Amount: amount, Type: "debit", Status: "completed"}
    if _, err := s.txRep.Create(ctx, tx, t); err != nil { return err }
    return tx.Commit()
}

func (s *TransactionService) Transfer(ctx context.Context, fromUserID, toUserID int64, amount float64) error {
    if amount <= 0 { return errors.New("amount_must_be_positive") }
    if fromUserID == toUserID { return errors.New("cannot_transfer_to_self") }

    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()

    if _, err := s.bal.AdjustBalance(ctx, tx, fromUserID, -amount); err != nil { return err }
    if _, err := s.bal.AdjustBalance(ctx, tx, toUserID,  amount); err != nil { return err }

    from := &fromUserID
    to := &toUserID
    t := &models.Transaction{FromUserID: from, ToUserID: to, Amount: amount, Type: "transfer", Status: "completed"}
    if _, err := s.txRep.Create(ctx, tx, t); err != nil { return err }
    return tx.Commit()
}

func (s *TransactionService) History(ctx context.Context, userID int64, limit, offset int) ([]models.Transaction, error) {
    return s.txRep.ListByUser(ctx, userID, limit, offset)
}

func (s *TransactionService) CurrentBalance(ctx context.Context, userID int64) (float64, error) {
    return s.bal.GetCurrent(ctx, userID)
}
