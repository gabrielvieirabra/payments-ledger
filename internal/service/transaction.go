package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gabrielvieirabra/payments-ledger/internal/domain"
	"github.com/gabrielvieirabra/payments-ledger/internal/repository"
	"github.com/gabrielvieirabra/payments-ledger/internal/worker"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrSameAccount         = errors.New("cannot transfer to the same account")
	ErrCurrencyMismatch    = errors.New("currency mismatch between accounts")
)

type TransactionService struct {
	accountRepo     *repository.AccountRepository
	entryRepo       *repository.EntryRepository
	transactionRepo *repository.TransactionRepository
	pool            *worker.Pool
}

func NewTransactionService(
	accountRepo *repository.AccountRepository,
	entryRepo *repository.EntryRepository,
	transactionRepo *repository.TransactionRepository,
	pool *worker.Pool,
) *TransactionService {
	return &TransactionService{
		accountRepo:     accountRepo,
		entryRepo:       entryRepo,
		transactionRepo: transactionRepo,
		pool:            pool,
	}
}

func (s *TransactionService) Transfer(ctx context.Context, req domain.CreateTransactionRequest) (domain.TransactionResult, error) {
	if req.FromAccountID == req.ToAccountID {
		return domain.TransactionResult{}, ErrSameAccount
	}

	var result domain.TransactionResult
	var execErr error

	errCh := make(chan error, 1)
	cmd := worker.Command{
		AccountID: req.FromAccountID,
		Exec: func(workerCtx context.Context) error {
			result, execErr = s.transferDirect(ctx, req)
			return execErr
		},
		Err: errCh,
	}

	if err := s.pool.Submit(cmd); err != nil {
		return domain.TransactionResult{}, fmt.Errorf("submit transfer command: %w", err)
	}

	if err := <-errCh; err != nil {
		return domain.TransactionResult{}, err
	}

	return result, nil
}

func (s *TransactionService) transferDirect(ctx context.Context, req domain.CreateTransactionRequest) (domain.TransactionResult, error) {
	fromAcc, err := s.accountRepo.GetByID(ctx, req.FromAccountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TransactionResult{}, fmt.Errorf("source %w", ErrAccountNotFound)
		}
		return domain.TransactionResult{}, err
	}

	toAcc, err := s.accountRepo.GetByID(ctx, req.ToAccountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TransactionResult{}, fmt.Errorf("destination %w", ErrAccountNotFound)
		}
		return domain.TransactionResult{}, err
	}

	if fromAcc.Currency != req.Currency || toAcc.Currency != req.Currency {
		return domain.TransactionResult{}, ErrCurrencyMismatch
	}

	pool := s.accountRepo.Pool()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return domain.TransactionResult{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			slog.Error("failed to rollback transaction", "error", rbErr)
		}
	}()

	// Lock accounts in consistent order to prevent deadlocks
	id1, id2 := req.FromAccountID, req.ToAccountID
	if id1.String() > id2.String() {
		id1, id2 = id2, id1
	}

	// Lock both accounts; use the locked row for balance check
	locked1, err := s.accountRepo.GetByIDForUpdate(ctx, tx, id1)
	if err != nil {
		return domain.TransactionResult{}, err
	}
	locked2, err := s.accountRepo.GetByIDForUpdate(ctx, tx, id2)
	if err != nil {
		return domain.TransactionResult{}, err
	}

	// Determine which locked row is the source account
	lockedFrom := locked1
	if id1 != req.FromAccountID {
		lockedFrom = locked2
	}

	if lockedFrom.Balance < req.Amount {
		return domain.TransactionResult{}, ErrInsufficientBalance
	}

	// Create transaction record
	txn, err := s.transactionRepo.Create(ctx, tx, req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		return domain.TransactionResult{}, err
	}

	// Create entries (debit from source, credit to destination)
	fromEntry, err := s.entryRepo.Create(ctx, tx, domain.CreateEntryParams{
		AccountID: req.FromAccountID,
		Amount:    -req.Amount,
	})
	if err != nil {
		return domain.TransactionResult{}, err
	}

	toEntry, err := s.entryRepo.Create(ctx, tx, domain.CreateEntryParams{
		AccountID: req.ToAccountID,
		Amount:    req.Amount,
	})
	if err != nil {
		return domain.TransactionResult{}, err
	}

	// Update balances
	updatedFrom, err := s.accountRepo.UpdateBalance(ctx, tx, req.FromAccountID, -req.Amount)
	if err != nil {
		return domain.TransactionResult{}, err
	}

	updatedTo, err := s.accountRepo.UpdateBalance(ctx, tx, req.ToAccountID, req.Amount)
	if err != nil {
		return domain.TransactionResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.TransactionResult{}, fmt.Errorf("commit transaction: %w", err)
	}

	return domain.TransactionResult{
		Transaction: txn,
		FromAccount: updatedFrom,
		ToAccount:   updatedTo,
		FromEntry:   fromEntry,
		ToEntry:     toEntry,
	}, nil
}

func (s *TransactionService) GetByID(ctx context.Context, id uuid.UUID) (domain.Transaction, error) {
	txn, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Transaction{}, ErrTransactionNotFound
		}
		return domain.Transaction{}, fmt.Errorf("get transaction: %w", err)
	}
	return txn, nil
}

func (s *TransactionService) ListByAccount(ctx context.Context, accountID uuid.UUID, limit, offset int32) ([]domain.Transaction, error) {
	return s.transactionRepo.ListByAccount(ctx, accountID, limit, offset)
}
