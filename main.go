package main

import (
	"log"
	"net/http"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// Note represents an note in our API
type Note struct {
	ID          uint   `gorm:"primaryKey" json:"id"` // Unique identifier for the note
	Name        string `json:"name"`                 // Name of the note
	Description string `json:"description"`          // Description of the note
}

var db *gorm.DB // Database connection

func main() {
	// Connect to SQLite database
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&Note{}); err != nil {
		log.Fatal(err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Define routes
	notesGroup := router.Group("/notes")
	{
		notesGroup.GET("", getNotes)          // Route to get all notes
		notesGroup.GET("/:id", getNote)       // Route to get a specific item by ID
		notesGroup.POST("", createNote)       // Route to create a new item
		notesGroup.PUT("/:id", updateNote)    // Route to update an existing item by ID
		notesGroup.DELETE("/:id", deleteNote) // Route to delete an item by ID
	}

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// Handlers

// getNotes handles the request to retrieve all notes
func getNotes(c *gin.Context) {
	var notes []Note
	if err := db.Find(&notes).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notes)
}

// getNote handles the request to retrieve a specific item by ID
func getNote(c *gin.Context) {
	var item Note
	id := c.Param("id")
	if err := db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// createNote handles the request to create a new item
func createNote(c *gin.Context) {
	var item Note
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.Name = strings.TrimSpace(item.Name)

	if item.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name cannot be blank"})
		return
	}

	if err := db.Create(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// updateNote handles the request to update an existing item by ID
func updateNote(c *gin.Context) {
	var item Note
	id := c.Param("id")
	if err := db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&item)
	c.JSON(http.StatusOK, item)
}

// deleteNote handles the request to delete an item by ID
func deleteNote(c *gin.Context) {
	var item Note
	id := c.Param("id")
	if err := db.Delete(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
