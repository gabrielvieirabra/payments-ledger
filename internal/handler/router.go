package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gabrielvieirabra/payments-ledger/internal/middleware"
	"github.com/gabrielvieirabra/payments-ledger/internal/repository"
	"github.com/gabrielvieirabra/payments-ledger/internal/service"
	"github.com/gabrielvieirabra/payments-ledger/internal/worker"
)

func NewRouter(pool *pgxpool.Pool, wp *worker.Pool) *gin.Engine {
	router := gin.New()
	_ = router.SetTrustedProxies(nil)
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.BodySizeLimit())

	accountRepo := repository.NewAccountRepository(pool)
	entryRepo := repository.NewEntryRepository(pool)
	transactionRepo := repository.NewTransactionRepository(pool)
	idempotencyRepo := repository.NewIdempotencyRepository(pool)

	accountSvc := service.NewAccountService(accountRepo)
	entrySvc := service.NewEntryService(entryRepo)
	transactionSvc := service.NewTransactionService(accountRepo, entryRepo, transactionRepo, wp)

	idempotencyMw := middleware.Idempotency(idempotencyRepo)

	healthH := NewHealthHandler(pool)
	accountH := NewAccountHandler(accountSvc)
	entryH := NewEntryHandler(entrySvc)
	transactionH := NewTransactionHandler(transactionSvc)

	router.GET("/healthz", healthH.Liveness)
	router.GET("/readyz", healthH.Readiness)

	v1 := router.Group("/api/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.POST("", idempotencyMw, accountH.Create)
			accounts.GET("", accountH.List)
			accounts.GET("/:id", accountH.GetByID)
			accounts.DELETE("/:id", accountH.Delete)
			accounts.GET("/:id/entries", entryH.ListByAccount)
			accounts.GET("/:id/transactions", transactionH.ListByAccount)
		}

		entries := v1.Group("/entries")
		{
			entries.GET("/:id", entryH.GetByID)
		}

		transactions := v1.Group("/transactions")
		{
			transactions.POST("", idempotencyMw, transactionH.Transfer)
			transactions.GET("/:id", transactionH.GetByID)
		}
	}

	return router
}
