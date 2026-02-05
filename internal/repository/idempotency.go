package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gabrielvieirabra/payments-ledger/internal/domain"
)

var ErrIdempotencyKeyNotFound = errors.New("idempotency key not found")

type IdempotencyRepository struct {
	pool *pgxpool.Pool
}

func NewIdempotencyRepository(pool *pgxpool.Pool) *IdempotencyRepository {
	return &IdempotencyRepository{pool: pool}
}

func (r *IdempotencyRepository) Find(ctx context.Context, key, method, path string) (domain.IdempotencyKey, error) {
	var ik domain.IdempotencyKey
	err := r.pool.QueryRow(ctx,
		`SELECT id, idempotency_key, method, path, status_code, response_body, created_at, expires_at
		 FROM idempotency_keys
		 WHERE idempotency_key = $1 AND method = $2 AND path = $3 AND expires_at > now()`,
		key, method, path,
	).Scan(&ik.ID, &ik.IdempotencyKey, &ik.Method, &ik.Path, &ik.StatusCode, &ik.ResponseBody, &ik.CreatedAt, &ik.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.IdempotencyKey{}, ErrIdempotencyKeyNotFound
		}
		return domain.IdempotencyKey{}, fmt.Errorf("find idempotency key: %w", err)
	}
	return ik, nil
}

func (r *IdempotencyRepository) Store(ctx context.Context, key, method, path string, statusCode int, responseBody []byte) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO idempotency_keys (idempotency_key, method, path, status_code, response_body)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (idempotency_key, method, path) DO NOTHING`,
		key, method, path, statusCode, responseBody,
	)
	if err != nil {
		return fmt.Errorf("store idempotency key: %w", err)
	}
	return nil
}
