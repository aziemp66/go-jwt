package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
	})

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success!!!",
		})
	})

	router.Run(":3000")
}
