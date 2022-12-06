package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("Secret Seklai")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string
	jwt.StandardClaims
}

var users = map[string]interface{}{
	"user1": "password1",
	"user2": "password2",
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Hour)
	claims["authorized"] = true
	claims["user"] = "username"

	tokenString, err := token.SignedString("blablabla")
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func main() {
	router := gin.Default()

	router.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
	})

	router.POST("/login", func(ctx *gin.Context) {
		var credentials Credentials

		ctx.ShouldBindJSON(&credentials)

		expectedPassword, ok := users[credentials.Username]

		if !ok || expectedPassword != credentials.Password {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 403,
			})
			return
		}

		expirationTime := time.Now().Add(time.Minute * 60)

		claims := &Claims{
			Username: credentials.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
			})
			return
		}

		ctx.SetCookie("token", tokenString, expirationTime.Second(), "/", "localhost", false, true)
	})

	router.Run(":3000")
}
