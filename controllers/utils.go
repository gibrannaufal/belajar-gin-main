package controllers

import (
	"net/http"
	"strconv"

	"github.com/gibrannaufal/belajar-main-gin/models"
	"github.com/gin-gonic/gin"
)

func CheckUserCredentials(c *gin.Context, userID int, username, password string) bool {
	// Validasi input
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password cannot be empty"})

		return false
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})

		return false
	}

	if user.Username != username || user.Password != password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})

		return false
	}

	return true
}

func StringToInt(s string) (int, error) {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func StringToInt64(s string) (int64, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func HandleOK(c *gin.Context, message string, data interface{}) bool {
	c.JSON(http.StatusOK, gin.H{"message": message, "data": data})
	return true
}

func HandleError(c *gin.Context, message string, err error, code int) bool {
	if err != nil {
		c.JSON(code, gin.H{
			"error":   message,
			"details": err.Error(),
		})
		return true
	}
	return false
}
