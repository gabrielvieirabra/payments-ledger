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

type EntryHandler struct {
	svc *service.EntryService
}

func NewEntryHandler(svc *service.EntryService) *EntryHandler {
	return &EntryHandler{svc: svc}
}

func (h *EntryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry id"})
		return
	}

	entry, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrEntryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return
		}
		slog.Error("failed to get entry", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get entry"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *EntryHandler) ListByAccount(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	var params domain.ListEntriesParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params.AccountID = accountID

	entries, err := h.svc.ListByAccount(c.Request.Context(), params)
	if err != nil {
		slog.Error("failed to list entries", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list entries"})
		return
	}

	c.JSON(http.StatusOK, entries)
}
