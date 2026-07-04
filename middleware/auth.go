package middleware

import (
	"strings"

	"ticket-system/utils"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "authorization header is required"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "expected format: Bearer <token>"})
			return
		}

		userID, err := utils.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
