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

	router.Static("/static", "docs")
	// router.GET("/static/*filepath", func(c *gin.Context) {
	// 	staticHandler := http.FileServer(http.Dir("docs"))
	// 	http.StripPrefix("/static/", staticHandler).ServeHTTP(c.Writer, c.Request)
	// })

	router.GET("/docs", func(c *gin.Context) {
		c.File("docs/index.html")
	})

	router.GET("api/manga/proxy", mangaController.MangabatProxy)
	router.GET("api/manga/source", mangaController.Source)
	router.GET("api/manga/search", mangaController.Search)
	router.GET("api/manga/index", mangaController.Index)
	router.GET("api/manga/chapter", mangaController.Chapter)
	router.GET("api/manga/detail", mangaController.Detail)
	router.GET("api/manga/image", mangaController.Image)
	router.ServeHTTP(w, r)
}
