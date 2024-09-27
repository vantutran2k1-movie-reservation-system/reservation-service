package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

type AuthMiddleware struct {
	userSessionRepo repositories.UserSessionRepository
}

func NewAuthMiddleware(userSessionRepo repositories.UserSessionRepository) *AuthMiddleware {
	return &AuthMiddleware{userSessionRepo: userSessionRepo}
}

func (m *AuthMiddleware) RequireAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenValue := utils.GetAuthorizationHeader(ctx.Request)
		if tokenValue == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		s, err := m.userSessionRepo.GetUserSession(m.userSessionRepo.GetUserSessionID(tokenValue))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		ctx.Set(constants.USER_SESSION, s)
		ctx.Next()
	}
}

func GetUserSession(ctx *gin.Context) (*models.UserSession, *errors.ApiError) {
	userSession, exist := ctx.Get(constants.USER_SESSION)
	if !exist {
		return nil, errors.InternalServerError("Can not get user id from request")
	}

	return userSession.(*models.UserSession), nil
}
