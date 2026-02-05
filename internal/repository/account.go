package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gabrielvieirabra/payments-ledger/internal/domain"
)

var ErrAccountHasReferences = errors.New("account has existing entries or transactions")

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{pool: pool}
}

func (r *AccountRepository) Create(ctx context.Context, req domain.CreateAccountRequest) (domain.Account, error) {
	var acc domain.Account
	err := r.pool.QueryRow(ctx,
		`INSERT INTO accounts (owner, currency) VALUES ($1, $2)
		 RETURNING id, owner, balance, currency, created_at, updated_at`,
		req.Owner, req.Currency,
	).Scan(&acc.ID, &acc.Owner, &acc.Balance, &acc.Currency, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return domain.Account{}, fmt.Errorf("create account: %w", err)
	}
	return acc, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Account, error) {
	var acc domain.Account
	err := r.pool.QueryRow(ctx,
		`SELECT id, owner, balance, currency, created_at, updated_at FROM accounts WHERE id = $1`,
		id,
	).Scan(&acc.ID, &acc.Owner, &acc.Balance, &acc.Currency, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return domain.Account{}, fmt.Errorf("get account: %w", err)
	}
	return acc, nil
}

func (r *AccountRepository) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (domain.Account, error) {
	var acc domain.Account
	err := tx.QueryRow(ctx,
		`SELECT id, owner, balance, currency, created_at, updated_at FROM accounts WHERE id = $1 FOR NO KEY UPDATE`,
		id,
	).Scan(&acc.ID, &acc.Owner, &acc.Balance, &acc.Currency, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return domain.Account{}, fmt.Errorf("get account for update: %w", err)
	}
	return acc, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, tx pgx.Tx, id uuid.UUID, amount int64) (domain.Account, error) {
	var acc domain.Account
	err := tx.QueryRow(ctx,
		`UPDATE accounts SET balance = balance + $1, updated_at = now() WHERE id = $2
		 RETURNING id, owner, balance, currency, created_at, updated_at`,
		amount, id,
	).Scan(&acc.ID, &acc.Owner, &acc.Balance, &acc.Currency, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return domain.Account{}, fmt.Errorf("update account balance: %w", err)
	}
	return acc, nil
}

func (r *AccountRepository) List(ctx context.Context, params domain.ListAccountsParams) ([]domain.Account, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, owner, balance, currency, created_at, updated_at FROM accounts
		 ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		params.Limit, params.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var acc domain.Account
		if err := rows.Scan(&acc.ID, &acc.Owner, &acc.Balance, &acc.Currency, &acc.CreatedAt, &acc.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, acc)
	}
	return accounts, rows.Err()
}

func (r *AccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM accounts WHERE id = $1`, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrAccountHasReferences
		}
		return fmt.Errorf("delete account: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *AccountRepository) Pool() *pgxpool.Pool {
	return r.pool
}
