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

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountHasReferences = errors.New("account has existing entries or transactions")
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) Create(ctx context.Context, req domain.CreateAccountRequest) (domain.Account, error) {
	return s.repo.Create(ctx, req)
}

func (s *AccountService) GetByID(ctx context.Context, id uuid.UUID) (domain.Account, error) {
	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Account{}, ErrAccountNotFound
		}
		return domain.Account{}, fmt.Errorf("get account: %w", err)
	}
	return acc, nil
}

func (s *AccountService) List(ctx context.Context, params domain.ListAccountsParams) ([]domain.Account, error) {
	return s.repo.List(ctx, params)
}

func (s *AccountService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccountNotFound
		}
		if errors.Is(err, repository.ErrAccountHasReferences) {
			return ErrAccountHasReferences
		}
		return fmt.Errorf("delete account: %w", err)
	}
	return nil
}
