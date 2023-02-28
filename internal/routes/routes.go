package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ivanbaug/go-eshops/internal/dbdriver"
)

var db *dbdriver.DB

func SetupRouter(rdb *dbdriver.DB) *gin.Engine {
	// Set db variable to the one passed in
	db = rdb

	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/add-product", AddProduct)

	r.GET("/get-products", GetProducts)

	return r
}
