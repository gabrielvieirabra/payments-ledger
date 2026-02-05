package middleware

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gabrielvieirabra/payments-ledger/internal/repository"
)

const IdempotencyKeyHeader = "Idempotency-Key"

type responseRecorder struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Idempotency(repo *repository.IdempotencyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader(IdempotencyKeyHeader)
		if key == "" {
			c.Next()
			return
		}

		if len(key) > 255 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "idempotency key must be at most 255 characters"})
			c.Abort()
			return
		}

		method := c.Request.Method
		path := c.FullPath()

		cached, err := repo.Find(c.Request.Context(), key, method, path)
		if err == nil {
			slog.Debug("idempotency cache hit",
				"key", key,
				"method", method,
				"path", path,
			)
			c.Data(cached.StatusCode, "application/json", cached.ResponseBody)
			c.Abort()
			return
		}
		if !errors.Is(err, repository.ErrIdempotencyKeyNotFound) {
			slog.Error("failed to check idempotency key", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			c.Abort()
			return
		}

		recorder := &responseRecorder{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = recorder

		c.Next()

		if c.IsAborted() {
			return
		}

		statusCode := c.Writer.Status()
		responseBody := recorder.body.Bytes()

		if statusCode < 200 || statusCode >= 300 {
			return
		}

		if err := repo.Store(c.Request.Context(), key, method, path, statusCode, responseBody); err != nil {
			slog.Error("failed to store idempotency key",
				"key", key,
				"error", err,
			)
		}
	}
}
