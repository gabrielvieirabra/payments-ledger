package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gabrielvieirabra/payments-ledger/internal/domain"
	"github.com/gabrielvieirabra/payments-ledger/internal/service"
)

type TransactionHandler struct {
	svc *service.TransactionService
}

func NewTransactionHandler(svc *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) Transfer(c *gin.Context) {
	var req domain.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.Transfer(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSameAccount):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrCurrencyMismatch):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrInsufficientBalance):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrAccountNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			slog.Error("failed to process transfer", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process transfer"})
		}
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	txn, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTransactionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		slog.Error("failed to get transaction", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get transaction"})
		return
	}

	c.JSON(http.StatusOK, txn)
}

func (h *TransactionHandler) ListByAccount(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	var params struct {
		Limit  int32 `form:"limit,default=10" binding:"min=1,max=100"`
		Offset int32 `form:"offset,default=0" binding:"min=0"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactions, err := h.svc.ListByAccount(c.Request.Context(), accountID, params.Limit, params.Offset)
	if err != nil {
		slog.Error("failed to list transactions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
