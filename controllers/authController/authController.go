package authController

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gibrannaufal/belajar-main-gin/models"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("29116d0494d43fd84d3474e6828005fcad889ead8ffd9d0c2b8f2b61920d4900a186206b26f628ffe039e48ef08fbdba312f7edd28a129c340a3936939cd9d106f0d18689c1e53f18587f0398c47bb03d67c45444f6c901bd2c7ddf63baa6b59f7f93e7ae26c338fd89c913ec98c758b918c8caca6520d4d97dc63249dd6ca303478935af5a294ee41cc8c7fb4be00b25dc72e2ae3a0adaf3e44d14186b030d92bbc3f74b7504d90240af2288a8ff01f60d60950a1e8830b54cfa70fd5b0729a74a9c5101c7bde679f1b67b882e74b38ff0d2b079c025dee67ffa380d57501a994d3a7c2559e172ecd82c6c788fc79eeb86afb822a22d6e30b4998244d43ce0f")

type Claims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundUser models.User
	result := models.DB.Where("username = ? AND password = ?", user.Username, user.Password).First(&foundUser)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: foundUser.Username,
		UserId:   foundUser.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": foundUser.Username, "token": tokenString})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
