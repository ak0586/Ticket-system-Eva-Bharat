package middleware

import (
	"strings"

	"ticket-system/utils"

	"github.com/gin-gonic/gin"
)

// Auth returns a Gin middleware function that intercepts incoming HTTP requests.
// Python equivalent: A FastAPI Depends() function or a starlette BaseHTTPMiddleware.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Grab the "Authorization" header from the incoming request.
		header := c.GetHeader("Authorization")
		if header == "" {
			// Abort stops the request from reaching the handler and returns a 401.
			c.AbortWithStatusJSON(401, gin.H{"error": "authorization header is required"})
			return
		}

		// 2. The header should look like "Bearer <token_string>". Split it into two parts.
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "expected format: Bearer <token>"})
			return
		}

		// 3. Pass the token string to our utility function to verify the signature and expiration.
		userID, err := utils.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired token"})
			return
		}

		// 4. Success! Store the user_id in the request context. 
		// The handlers (like CreateTicket) will retrieve this user_id later.
		c.Set("user_id", userID)
		
		// 5. c.Next() tells Gin to proceed to the actual handler function.
		c.Next()
	}
}
