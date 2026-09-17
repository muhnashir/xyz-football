package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"github.com/xyz-corp/xyz-football-api/pkg/response"
	"go.uber.org/zap"
)

// Recovery converts panics into a generic 500 response so internal errors are never leaked (NFR Keamanan: Info leak).
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", zap.Any("error", r), zap.String("path", c.Request.URL.Path))
				response.Error(c, apperror.Internal("Terjadi kesalahan pada server"))
				c.Abort()
			}
		}()
		c.Next()
	}
}
