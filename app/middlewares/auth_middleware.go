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
	featureFlagRepo repositories.FeatureFlagRepository
}

func NewAuthMiddleware(userSessionRepo repositories.UserSessionRepository, featureFlagRepo repositories.FeatureFlagRepository) *AuthMiddleware {
	return &AuthMiddleware{userSessionRepo: userSessionRepo, featureFlagRepo: featureFlagRepo}
}

func (m *AuthMiddleware) RequireAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenValue := utils.GetAuthorizationHeader(ctx.Request)
		if tokenValue == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		s, err := m.userSessionRepo.GetUserSession(m.userSessionRepo.GetUserSessionID(tokenValue))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if s == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		ctx.Set(constants.UserSession, s)
		ctx.Next()
	}
}

func (m *AuthMiddleware) RequireFeatureFlagMiddleware(flagName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := GetUserSession(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(err.StatusCode, gin.H{"error": "can not get feature flags of user"})
			return
		}

		hasFlagEnabled := m.featureFlagRepo.HasFlagEnabled(session.Email, flagName)
		if !hasFlagEnabled {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission error"})
			return
		}

		ctx.Next()
	}
}

func GetUserSession(ctx *gin.Context) (*models.UserSession, *errors.ApiError) {
	userSession, exist := ctx.Get(constants.UserSession)
	if !exist {
		return nil, errors.InternalServerError("Can not get user id from request")
	}

	return userSession.(*models.UserSession), nil
}
