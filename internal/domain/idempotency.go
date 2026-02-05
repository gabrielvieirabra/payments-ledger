package domain

import (
	"time"

	"github.com/google/uuid"
)

type IdempotencyKey struct {
	ID             uuid.UUID `json:"id"`
	IdempotencyKey string    `json:"idempotency_key"`
	Method         string    `json:"method"`
	Path           string    `json:"path"`
	StatusCode     int       `json:"status_code"`
	ResponseBody   []byte    `json:"response_body"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}
