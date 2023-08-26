package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"little_mangamee/controller"
	log "little_mangamee/logger"
	"little_mangamee/middleware"
	"little_mangamee/service"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	log.InitZerolog()

	router := gin.New()
	router.Use(middleware.Logger(log.Logger))
	router.Use(cors.Default())
	router.Use(gin.Recovery())

	mangaService := service.NewMangaService()
	mangaController := controller.NewMangaController(mangaService)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Message": "Hello, Little mangamee!",
			"version": "1.0.0",
		})
	})

	router.GET("api/manga/index", mangaController.Index)
	router.ServeHTTP(w, r)
}
