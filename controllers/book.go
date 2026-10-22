package controllers

import (
	"log"
	"net/http"
	"tac/gonic1/dal"

	"github.com/gin-gonic/gin"
)

type bookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

// GET /books
// Find all books
func FindBooks(c *gin.Context) {
	var books []dal.Book
	dal.DB.Find(&books)
	log.Printf("[INFO] get all books #%d", len(books))
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// GET /books/:id
// Find a book
func FindBook(c *gin.Context) {
	// Get model if exist
	var book dal.Book
	id := c.Param("id")
	if err := dal.DB.Where("id = ?", id).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record " + id + " not found!"})
		return
	}
	log.Printf("[INFO] get book %s", id)
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// POST /books
// Create new book
func CreateBook(c *gin.Context) {
	// Validate input
	var input bookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create book
	book := dal.Book{Title: input.Title, Author: input.Author}
	dal.DB.Create(&book)
	log.Printf("[INFO] created book '%v'", book)
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// PUT /books/:id
// Update a book
func UpdateBook(c *gin.Context) {
	// Get model if exist
	var book dal.Book
	id := c.Param("id")
	if err := dal.DB.Where("id = ?", id).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record " + id + "not found!"})
		return
	}

	// Validate input
	var input bookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book.Title = input.Title
	book.Author = input.Author
	log.Printf("[INFO] updating book '%v'", book)
	dal.DB.Save(&book)

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// DELETE /books/:id
// Delete a book
func DeleteBook(c *gin.Context) {
	// Get model if exist
	var book dal.Book
	if err := dal.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	log.Printf("[INFO] deleting book '%v'", book)
	dal.DB.Delete(&book)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
