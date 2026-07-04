package handlers

import (
	"ticket-system/db"
	"ticket-system/models"

	"github.com/gin-gonic/gin"
)

type createTicketInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func CreateTicket(c *gin.Context) {
	var input createTicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	ticket := models.Ticket{
		Title:       input.Title,
		Description: input.Description,
		Status:      "open",
		UserID:      userID,
	}

	if err := db.DB.Create(&ticket).Error; err != nil {
		c.JSON(500, gin.H{"error": "could not create ticket"})
		return
	}

	c.JSON(201, ticket)
}

func ListTickets(c *gin.Context) {
	userID := c.GetUint("user_id")

	var tickets []models.Ticket
	db.DB.Where("user_id = ?", userID).Find(&tickets)

	c.JSON(200, tickets)
}

func GetTicket(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var ticket models.Ticket
	if err := db.DB.Where("id = ? AND user_id = ?", id, userID).First(&ticket).Error; err != nil {
		c.JSON(404, gin.H{"error": "ticket not found"})
		return
	}

	c.JSON(200, ticket)
}

// validTransitions defines the only allowed status moves
var validTransitions = map[string]string{
	"open":        "in_progress",
	"in_progress": "closed",
}

type updateStatusInput struct {
	Status string `json:"status" binding:"required"`
}

func UpdateStatus(c *gin.Context) {
	var input updateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	id := c.Param("id")

	var ticket models.Ticket
	if err := db.DB.Where("id = ? AND user_id = ?", id, userID).First(&ticket).Error; err != nil {
		c.JSON(404, gin.H{"error": "ticket not found"})
		return
	}

	if ticket.Status == "closed" {
		c.JSON(400, gin.H{"error": "closed ticket cannot be reopened"})
		return
	}

	allowed, ok := validTransitions[ticket.Status]
	if !ok || allowed != input.Status {
		c.JSON(400, gin.H{
			"error": "invalid status transition",
			"hint":  "allowed: open → in_progress → closed",
		})
		return
	}

	db.DB.Model(&ticket).Update("status", input.Status)
	ticket.Status = input.Status

	c.JSON(200, ticket)
}

