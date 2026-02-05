package domain

import (
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateEntryParams struct {
	AccountID uuid.UUID `json:"account_id"`
	Amount    int64     `json:"amount"`
}

type ListEntriesParams struct {
	AccountID uuid.UUID `form:"account_id" binding:"required"`
	Limit     int32     `form:"limit,default=10" binding:"min=1,max=100"`
	Offset    int32     `form:"offset,default=0" binding:"min=0"`
}
