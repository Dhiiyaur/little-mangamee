package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"little_mangamee/controller"
	"little_mangamee/docs"
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

	router.GET("/static/mangamee_collection.yml", func(c *gin.Context) {
		Colection, err := docs.Colection.ReadFile("mangamee_collection.yml")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading template")
			return
		}
		c.Data(http.StatusOK, "text/html", Colection)
	})

	router.GET("/static/redoc.standalone.js", func(c *gin.Context) {
		Colection, err := docs.RedocJS.ReadFile("redoc.standalone.js")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading template")
			return
		}
		c.Data(http.StatusOK, "text/html", Colection)
	})

	router.GET("/docs/*any", func(c *gin.Context) {
		htmlTemplate, err := docs.HtmlBase.ReadFile("index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading template")
			return
		}
		c.Data(http.StatusOK, "text/html", htmlTemplate)
	})

	router.GET("api/manga/proxy", mangaController.MangabatProxy)
	router.GET("api/manga/redirect", mangaController.Redirect)
	router.GET("api/manga/source", mangaController.Source)
	router.GET("api/manga/search", mangaController.Search)
	router.GET("api/manga/index", mangaController.Index)
	router.GET("api/manga/chapter", mangaController.Chapter)
	router.GET("api/manga/detail", mangaController.Detail)
	router.GET("api/manga/image", mangaController.Image)
	router.ServeHTTP(w, r)
}
