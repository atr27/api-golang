package auth

import (
	"auth/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
)

const (
	USER = "admin"
	PASSWORD = "admin123"
	SECRET = "secret"
)

func LoginHandler(c *gin.Context) {
	var user models.Credential
	arr := c.BindJSON(&user)
	if arr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	if user.Username != USER {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user invalid"})
		return
	} else {
		if user.Password != PASSWORD {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "password invalid"})
			return
		}
	}

	claim := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		Issuer: "admin",
		IssuedAt: time.Now().Unix(),
	}

	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := sign.SignedString([]byte(SECRET))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		c.Abort()
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "message": "success login"})	
}