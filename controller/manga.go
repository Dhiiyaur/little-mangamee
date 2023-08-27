package controller

import "github.com/gin-gonic/gin"

type MangaController interface {
	Index(c *gin.Context)
	Source(c *gin.Context)
	Chapter(c *gin.Context)
	Detail(c *gin.Context)
	Image(c *gin.Context)
	Search(c *gin.Context)
	MangabatProxy(c *gin.Context)
}
