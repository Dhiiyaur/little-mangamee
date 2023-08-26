package controller

import (
	"fmt"
	"little_mangamee/response"
	"little_mangamee/service"
	"little_mangamee/utils"

	"github.com/gin-gonic/gin"
)

type mangaControllerImpl struct {
	service service.MangaService
}

func NewMangaController(service service.MangaService) MangaController {
	return &mangaControllerImpl{
		service: service,
	}
}

func (s *mangaControllerImpl) Source(c *gin.Context) {

}

func (s *mangaControllerImpl) Index(c *gin.Context) {

	//
	ctx := c.Request.Context()
	source := c.Query("source")
	pageNumber := c.Query("page")

	switch source {
	case "mangabat":

		data, err := s.service.MangabatIndex(ctx, pageNumber)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangaread":

		data, err := s.service.MangareadIndex(ctx, pageNumber)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangatown":

		data, err := s.service.MangatownIndex(ctx, pageNumber)
		if err != nil {
			fmt.Println(err)
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "maidmy":

		data, err := s.service.MaidmyIndex(ctx, pageNumber)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	default:
		response.ErrorResponse(c, utils.ERR_BAD_REQUEST, utils.ERR_BAD_REQUEST.Error())
		return
	}

}
