package controller

import "github.com/gin-gonic/gin"

type MangaController interface {
	Index(c *gin.Context)
	Source(c *gin.Context)
}
