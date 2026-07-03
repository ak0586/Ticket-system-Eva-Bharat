package handlers

import (
	"ticket-system/db"
	"ticket-system/models"
	"ticket-system/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registerInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var input registerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// check if this email is already taken
	var existing models.User
	if result := db.DB.Where("email = ?", input.Email).First(&existing); result.Error == nil {
		c.JSON(409, gin.H{"error": "email already registered"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "something went wrong"})
		return
	}

	user := models.User{
		Email:    input.Email,
		Password: string(hash),
	}
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "could not create user"})
		return
	}

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

	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// don't say "user not found" — just say invalid credentials
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
