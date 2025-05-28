package middleware

import (
	"pos-master/services/authservices"
	"pos-master/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			utils.RespondWithError(c, 401, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := authservices.ValidateToken(tokenString)
		if err != nil {
			utils.RespondWithError(c, 401, "Invalid Token")
			c.Abort()
			return
		}

		c.Set("userID", claims.Subject)

		c.Next()
	}
}
