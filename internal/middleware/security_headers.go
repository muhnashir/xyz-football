package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders sets common hardening headers (NFR Keamanan: Header).
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Next()
	}
}
