package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gabrielvieirabra/payments-ledger/internal/domain"
	"github.com/gabrielvieirabra/payments-ledger/internal/repository"
)

var ErrEntryNotFound = errors.New("entry not found")

type EntryService struct {
	repo *repository.EntryRepository
}

func NewEntryService(repo *repository.EntryRepository) *EntryService {
	return &EntryService{repo: repo}
}

func (s *EntryService) GetByID(ctx context.Context, id uuid.UUID) (domain.Entry, error) {
	entry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Entry{}, ErrEntryNotFound
		}
		return domain.Entry{}, fmt.Errorf("get entry: %w", err)
	}
	return entry, nil
}

func (s *EntryService) ListByAccount(ctx context.Context, params domain.ListEntriesParams) ([]domain.Entry, error) {
	return s.repo.ListByAccount(ctx, params)
}
