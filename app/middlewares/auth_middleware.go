package middlewares

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
)

const USER_ID = "user_id"

type AuthMiddleware struct {
	userSessionRepo repositories.UserSessionRepository
}

func NewAuthMiddleware(userSessionRepo repositories.UserSessionRepository) *AuthMiddleware {
	return &AuthMiddleware{userSessionRepo: userSessionRepo}
}

func (m *AuthMiddleware) RequireJwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenValue := auth.GetAuthTokenFromRequest(ctx.Request)
		if tokenValue == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		claims := &auth.JwtClaims{}
		token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (any, error) {
			return auth.JwtKey, nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		s, err := m.userSessionRepo.GetUserSession(m.userSessionRepo.GetUserSessionID(tokenValue))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		ctx.Set(USER_ID, s.UserID)
		ctx.Next()
	}
}

func (m *AuthMiddleware) RequireBasicAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenValue := auth.GetAuthTokenFromRequest(ctx.Request)
		if tokenValue == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		s, err := m.userSessionRepo.GetUserSession(m.userSessionRepo.GetUserSessionID(tokenValue))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		ctx.Set(USER_ID, s.UserID)
		ctx.Next()
	}
}

func GetUserID(ctx *gin.Context) (uuid.UUID, *errors.ApiError) {
	userID, exist := ctx.Get(USER_ID)
	if !exist {
		return uuid.Nil, errors.InternalServerError("Can not get user id from request")
	}

	return userID.(uuid.UUID), nil
}
