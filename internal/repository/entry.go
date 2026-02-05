package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gabrielvieirabra/payments-ledger/internal/domain"
)

type EntryRepository struct {
	pool *pgxpool.Pool
}

func NewEntryRepository(pool *pgxpool.Pool) *EntryRepository {
	return &EntryRepository{pool: pool}
}

func (r *EntryRepository) Create(ctx context.Context, tx pgx.Tx, params domain.CreateEntryParams) (domain.Entry, error) {
	var entry domain.Entry
	err := tx.QueryRow(ctx,
		`INSERT INTO entries (account_id, amount) VALUES ($1, $2)
		 RETURNING id, account_id, amount, created_at`,
		params.AccountID, params.Amount,
	).Scan(&entry.ID, &entry.AccountID, &entry.Amount, &entry.CreatedAt)
	if err != nil {
		return domain.Entry{}, fmt.Errorf("create entry: %w", err)
	}
	return entry, nil
}

func (r *EntryRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Entry, error) {
	var entry domain.Entry
	err := r.pool.QueryRow(ctx,
		`SELECT id, account_id, amount, created_at FROM entries WHERE id = $1`,
		id,
	).Scan(&entry.ID, &entry.AccountID, &entry.Amount, &entry.CreatedAt)
	if err != nil {
		return domain.Entry{}, fmt.Errorf("get entry: %w", err)
	}
	return entry, nil
}

func (r *EntryRepository) ListByAccount(ctx context.Context, params domain.ListEntriesParams) ([]domain.Entry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, account_id, amount, created_at FROM entries
		 WHERE account_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		params.AccountID, params.Limit, params.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list entries: %w", err)
	}
	defer rows.Close()

	var entries []domain.Entry
	for rows.Next() {
		var entry domain.Entry
		if err := rows.Scan(&entry.ID, &entry.AccountID, &entry.Amount, &entry.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan entry: %w", err)
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}
