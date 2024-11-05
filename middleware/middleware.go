package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	SECRET = "secret"
)

func AuthValidator(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "token required"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _,invalid := token.Method.(*jwt.SigningMethodHMAC); !invalid {
			return nil, fmt.Errorf("Invalid token", token.Header["alg"])
		}
		return []byte(SECRET), nil
	})

	if token != nil && err == nil {
		fmt.Println("token verified")
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "token unauthorized", "error": err.Error()})
		c.Abort()
		return
	}
	
}
