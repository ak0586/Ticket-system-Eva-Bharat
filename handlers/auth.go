package handlers

import (
	"ticket-system/db"
	"ticket-system/models"
	"ticket-system/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// registerInput acts exactly like a Pydantic BaseModel in FastAPI.
// The `binding:"required,email"` tags tell Gin to automatically validate the input.
type registerInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var input registerInput
	
	// ShouldBindJSON reads the request body and validates it against the struct tags.
	// If the user sends an invalid email, this will throw an error automatically.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if this email is already taken to prevent duplicates.
	var existing models.User
	if result := db.DB.Where("email = ?", input.Email).First(&existing); result.Error == nil {
		c.JSON(409, gin.H{"error": "email already registered"})
		return
	}

	// Hash the password using bcrypt. NEVER store plain-text passwords.
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "something went wrong"})
		return
	}

	// Prepare the User model for insertion into the database.
	user := models.User{
		Email:    input.Email,
		Password: string(hash),
	}
	
	// Execute the INSERT query.
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "could not create user"})
		return
	}

	// Return 201 Created on success.
	c.JSON(201, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

type loginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input loginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Query the database for the user by email.
	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// Security Best Practice: Don't tell the user "email not found". 
		// It helps hackers enumerate valid accounts. Just say "invalid credentials".
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// Compare the hash stored in the DB with the plain-text password the user typed.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate a JWT since they successfully proved who they are.
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	// Return the token so the client can use it in future requests.
	c.JSON(200, gin.H{"token": token})
}
