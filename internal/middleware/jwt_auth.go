package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"github.com/xyz-corp/xyz-football-api/pkg/jwt"
	"github.com/xyz-corp/xyz-football-api/pkg/response"
)

const ContextUserUUIDKey = "user_uuid"
const ContextUserEmailKey = "user_email"

func JWTAuth(manager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, apperror.Unauthorized("Token tidak ditemukan"))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		claims, err := manager.Parse(tokenString)
		if err != nil {
			response.Error(c, apperror.Unauthorized("Token tidak valid atau kedaluwarsa"))
			c.Abort()
			return
		}

		c.Set(ContextUserUUIDKey, claims.Sub.String())
		c.Set(ContextUserEmailKey, claims.Email)
		c.Next()
	}
}
