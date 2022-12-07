package main

import (
	"fmt"
	"net/http"
	"strings"
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
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
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
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
	})

	router.GET("/home", func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
		}
		tokenString := strings.Split(authHeader, " ")[1]

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			fmt.Println("Err 2")
			if err == jwt.ErrSignatureInvalid {
				ctx.Writer.WriteHeader(http.StatusUnauthorized) //signing method invalid
				return
			}
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			fmt.Println("Err 3")
			ctx.Writer.WriteHeader(http.StatusUnauthorized) //secret key invalid
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Hello %s", claims.Username),
		})
	})

	router.Run(":3000")
}
