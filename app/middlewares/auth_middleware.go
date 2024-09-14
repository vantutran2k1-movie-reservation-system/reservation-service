package middlewares

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := auth.GetAuthTokenFromRequest(c.Request)
		if tokenValue == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		claims := &auth.JwtClaims{}
		token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (any, error) {
			return auth.JwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		s, err := sessions.GetSession(tokenValue)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("user_id", s.UserID)
		c.Next()
	}
}
