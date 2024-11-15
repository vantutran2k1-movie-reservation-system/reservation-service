package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/context"
)

type ContextMiddleware struct{}

func NewContextMiddleware() *ContextMiddleware {
	return &ContextMiddleware{}
}

func (m *ContextMiddleware) AddRequestContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := context.RequestContext{RequestID: uuid.New()}
		context.SetRequestContext(ctx, c)
		ctx.Next()
	}
}
