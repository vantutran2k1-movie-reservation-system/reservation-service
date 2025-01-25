package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"net/http"
)

type RateLimitMiddleware struct{}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{}
}

func (m *RateLimitMiddleware) NotExceedMaxRequests(limiter services.RateLimiterService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		allowed, retryAfter := limiter.Allow(ctx.ClientIP())
		if !allowed {
			ctx.Header(constants.RetryAfter, fmt.Sprintf("%.0f", retryAfter.Seconds()))
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": retryAfter.Seconds(),
			})
			return
		}

		ctx.Next()
	}
}
