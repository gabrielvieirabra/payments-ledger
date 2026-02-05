package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const maxBodySize = 1 << 20 // 1 MB

func BodySizeLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBodySize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			c.Abort()
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodySize)
		c.Next()
	}
}
