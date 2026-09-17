package middleware

import (
	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"github.com/xyz-corp/xyz-football-api/pkg/response"
)

// RateLimiter returns Gin middleware enforcing the given rate (e.g. "5-M" = 5 requests/minute),
// used to mitigate brute force login attempts (NFR Keamanan).
func RateLimiter(rate string) gin.HandlerFunc {
	formatted, err := limiter.NewRateFromFormatted(rate)
	if err != nil {
		formatted, _ = limiter.NewRateFromFormatted("5-M")
	}
	store := memory.NewStore()
	instance := limiter.New(store, formatted)

	mw := ginlimiter.NewMiddleware(instance, ginlimiter.WithErrorHandler(func(c *gin.Context, e error) {
		response.Error(c, apperror.Internal("rate limiter error"))
		c.Abort()
	}), ginlimiter.WithLimitReachedHandler(func(c *gin.Context) {
		response.Error(c, apperror.TooManyRequests("Terlalu banyak percobaan, coba lagi nanti"))
		c.Abort()
	}))

	return mw
}
