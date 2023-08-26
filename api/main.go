package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/api/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	r.POST("/api/post", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, post!")
	})

	r.GET("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Little mangamee!")
	})

	r.Run()
}
