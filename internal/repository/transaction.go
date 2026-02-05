package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gabrielvieirabra/payments-ledger/internal/domain"
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{pool: pool}
}

func (r *TransactionRepository) Create(ctx context.Context, tx pgx.Tx, fromID, toID uuid.UUID, amount int64) (domain.Transaction, error) {
	var txn domain.Transaction
	err := tx.QueryRow(ctx,
		`INSERT INTO transactions (from_account_id, to_account_id, amount) VALUES ($1, $2, $3)
		 RETURNING id, from_account_id, to_account_id, amount, created_at`,
		fromID, toID, amount,
	).Scan(&txn.ID, &txn.FromAccountID, &txn.ToAccountID, &txn.Amount, &txn.CreatedAt)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("create transaction: %w", err)
	}
	return txn, nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Transaction, error) {
	var txn domain.Transaction
	err := r.pool.QueryRow(ctx,
		`SELECT id, from_account_id, to_account_id, amount, created_at FROM transactions WHERE id = $1`,
		id,
	).Scan(&txn.ID, &txn.FromAccountID, &txn.ToAccountID, &txn.Amount, &txn.CreatedAt)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("get transaction: %w", err)
	}
	return txn, nil
}

func (r *TransactionRepository) ListByAccount(ctx context.Context, accountID uuid.UUID, limit, offset int32) ([]domain.Transaction, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, from_account_id, to_account_id, amount, created_at FROM transactions
		 WHERE from_account_id = $1 OR to_account_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		accountID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var txn domain.Transaction
		if err := rows.Scan(&txn.ID, &txn.FromAccountID, &txn.ToAccountID, &txn.Amount, &txn.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		transactions = append(transactions, txn)
	}
	return transactions, rows.Err()
}
