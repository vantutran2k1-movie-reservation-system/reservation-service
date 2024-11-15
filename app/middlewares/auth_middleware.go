package middlewares

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/context"
	"net/http"

	"github.com/gin-gonic/gin"
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

		reqContext, apiErr := context.GetRequestContext(ctx)
		if apiErr != nil {
			ctx.AbortWithStatusJSON(apiErr.StatusCode, gin.H{"error": apiErr.Error()})
			return
		}
		reqContext.UserSession = s

		context.SetRequestContext(ctx, reqContext)
		ctx.Next()
	}
}

func (m *AuthMiddleware) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenValue := utils.GetAuthorizationHeader(ctx.Request)
		if tokenValue == "" {
			ctx.Next()
			return
		}

		s, err := m.userSessionRepo.GetUserSession(m.userSessionRepo.GetUserSessionID(tokenValue))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if s != nil {
			reqContext, apiErr := context.GetRequestContext(ctx)
			if apiErr != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": apiErr.Error()})
				return
			}
			reqContext.UserSession = s

			context.SetRequestContext(ctx, reqContext)
		}

		ctx.Next()
	}
}

func (m *AuthMiddleware) RequireFeatureFlagMiddleware(flagName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqContext, err := context.GetRequestContext(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(err.StatusCode, gin.H{"error": err.Error()})
			return
		}
		if reqContext.UserSession == nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "can not get feature flags of user"})
			return
		}

		hasFlagEnabled := m.featureFlagRepo.HasFlagEnabled(reqContext.UserSession.Email, flagName)
		if !hasFlagEnabled {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission error"})
			return
		}

		ctx.Next()
	}
}
