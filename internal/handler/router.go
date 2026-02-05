package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gabrielvieirabra/payments-ledger/internal/repository"
	"github.com/gabrielvieirabra/payments-ledger/internal/service"
)

func NewRouter(pool *pgxpool.Pool) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	accountRepo := repository.NewAccountRepository(pool)
	entryRepo := repository.NewEntryRepository(pool)
	transactionRepo := repository.NewTransactionRepository(pool)

	accountSvc := service.NewAccountService(accountRepo)
	entrySvc := service.NewEntryService(entryRepo)
	transactionSvc := service.NewTransactionService(accountRepo, entryRepo, transactionRepo)

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
			accounts.POST("", accountH.Create)
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
			transactions.POST("", transactionH.Transfer)
			transactions.GET("/:id", transactionH.GetByID)
		}
	}

	return router
}
