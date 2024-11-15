package context

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
)

const contextKey = "contextKey"

type RequestContext struct {
	RequestID   uuid.UUID
	UserSession *models.UserSession
}

func GetRequestContext(ctx *gin.Context) (RequestContext, *errors.ApiError) {
	c, exist := ctx.Get(contextKey)
	if !exist || c == nil {
		return RequestContext{}, errors.InternalServerError("can not get request context")
	}

	return c.(RequestContext), nil
}

func SetRequestContext(ctx *gin.Context, c RequestContext) {
	ctx.Set(contextKey, c)
}
