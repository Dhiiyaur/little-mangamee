package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func VercelHandler(w http.ResponseWriter, r *http.Request) {
	// Create a Gin router
	router := gin.New()
	router.Use(cors.Default())
	router.Use(gin.Recovery())

	router.GET("/api/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.POST("/api/post", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, post!")
	})

	router.GET("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Little mangamee!")
	})

	// Define your routes here using router.GET(), router.POST(), etc.

	// Use Gin's ServeHTTP method to handle the request
	router.ServeHTTP(w, r)
}

// func main() {
// 	r := gin.Default()
// 	r.Use(cors.Default())

// 	r.GET("/api/hello", func(c *gin.Context) {
// 		c.String(http.StatusOK, "Hello, World!")
// 	})

// 	r.POST("/api/post", func(c *gin.Context) {
// 		c.String(http.StatusOK, "Hello, post!")
// 	})

// 	r.GET("/api", func(c *gin.Context) {
// 		c.String(http.StatusOK, "Hello, Little mangamee!")
// 	})

// 	r.Run()
// }
