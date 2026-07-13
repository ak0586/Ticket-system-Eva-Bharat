package handlers

import "github.com/gin-gonic/gin"

// HealthCheck is a simple endpoint used by load balancers and Docker to see if the API is alive.
// Python equivalent: 
// @app.get("/health")
// def health_check(): return {"status": "ok"}
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
