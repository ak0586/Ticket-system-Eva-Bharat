package main

import (
	"log"
	"os"

	"ticket-system/db"
	"ticket-system/handlers"
	"ticket-system/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables from the .env file.
	// Python equivalent: load_dotenv() from python-dotenv.
	godotenv.Load()

	// 2. Initialize the database connection and auto-migrate schemas.
	// Python equivalent: db.init_db() or metadata.create_all(bind=engine)
	db.Init()

	// 3. Create a new Gin router with default logging and recovery middleware.
	// Python equivalent: app = FastAPI()
	router := gin.Default()

	// 4. Define a simple health check route (GET /health).
	// Python equivalent: @app.get("/health")
	router.GET("/health", handlers.HealthCheck)

	// 5. Create a route group for authentication endpoints.
	// Python equivalent: router = APIRouter(prefix="/auth")
	auth := router.Group("/auth")
	{
		// Map the Register handler to POST /auth/register
		auth.POST("/register", handlers.Register)
		// Map the Login handler to POST /auth/login
		auth.POST("/login", handlers.Login)
	}

	// 6. Create a route group for ticket endpoints.
	tickets := router.Group("/tickets")
	// 7. Apply the authentication middleware to the ENTIRE /tickets group.
	// Any route inside this block requires a valid JWT token.
	// Python equivalent: router = APIRouter(prefix="/tickets", dependencies=[Depends(verify_token)])
	tickets.Use(middleware.Auth())
	{
		tickets.POST("", handlers.CreateTicket)           // POST /tickets
		tickets.GET("", handlers.ListTickets)             // GET /tickets
		tickets.GET("/:id", handlers.GetTicket)           // GET /tickets/1
		tickets.PATCH("/:id/status", handlers.UpdateStatus) // PATCH /tickets/1/status
	}

	// 8. Get the PORT from the environment, default to 8080 if not found.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 9. Start the HTTP server.
	// Python equivalent: uvicorn.run(app, host="0.0.0.0", port=8080)
	log.Printf("starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
