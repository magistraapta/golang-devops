package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


func main() {
	router := gin.Default()

	router.LoadHTMLFiles("index.html")

	router.GET("/", RenderIndex)
	router.GET("/health", CheckHealth)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	fmt.Println("Server is running at port: 8080")
	if err := router.Run(":8080"); err != nil {
		fmt.Printf("Server is not running: %v", err)
	}
}

func RenderIndex(ctx *gin.Context) {
	ctx.HTML(200, "index.html", nil)
}

func CheckHealth(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "service is running",
	})
}