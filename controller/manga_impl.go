package controller

import (
	"fmt"
	"io/ioutil"
	"little_mangamee/entity"
	log "little_mangamee/logger"
	"little_mangamee/response"
	"little_mangamee/service"
	"little_mangamee/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type mangaControllerImpl struct {
	service service.MangaService
}

var mangaseeJsonData []entity.IndexData
var mangaseeJsonDataExpired time.Time

func NewMangaController(service service.MangaService) MangaController {
	return &mangaControllerImpl{
		service: service,
	}
}

func (s *mangaControllerImpl) Source(c *gin.Context) {

	data := []string{"mangabat", "mangaread", "mangatown", "maidmy", "asuracomic", "manganato", "manganelo", "mangasee"}
	response.SuccesResponse(c, data)
	return
}

func (s *mangaControllerImpl) Index(c *gin.Context) {

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

	case "asuracomic":
		data, err := s.service.AsuraComicIndex(ctx, pageNumber)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "manganato":
		data, err := s.service.ManganatoIndex(ctx, pageNumber)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "manganelo":
		data, err := s.service.ManganeloIndex(ctx, pageNumber)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangasee":

		var err error
		if len(mangaseeJsonData) == 0 {
			mangaseeJsonData, err = s.service.MangaseeIndex(ctx)
			if err != nil {
				log.Info().Err(err).Msg("err")
				response.ErrorResponse(c, err, nil)
				return
			}
			mangaseeJsonDataExpired = time.Now()
		}

		if time.Now().UTC().After(mangaseeJsonDataExpired.Add(10 * time.Minute)) {
			log.Info().Msg("Expired Data, Refetching")
			mangaseeJsonData, err = s.service.MangaseeIndex(ctx)
			if err != nil {
				log.Info().Err(err).Msg("err")
				response.ErrorResponse(c, err, nil)
				return
			}
			mangaseeJsonDataExpired = time.Now()
		}

		pageNumberInt, err := strconv.Atoi(pageNumber)
		if err != nil {
			log.Info().Err(err).Msg("err")
			response.ErrorResponse(c, err, nil)
			return
		}

		data, err := utils.PaginateIndex(mangaseeJsonData, pageNumberInt, 15)
		if err != nil {
			log.Info().Err(err).Msg("err")
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

func (s *mangaControllerImpl) Chapter(c *gin.Context) {

	ctx := c.Request.Context()
	source := c.Query("source")
	mangaId := c.Query("mangaid")

	switch source {
	case "mangabat":

		data, err := s.service.MangabatChapter(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangaread":

		data, err := s.service.MangareadChapter(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangatown":

		data, err := s.service.MangatownChapter(ctx, mangaId)
		if err != nil {
			fmt.Println(err)
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "maidmy":

		data, err := s.service.MaidmyChapter(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "asuracomic":
		data, err := s.service.AsuraComicChapter(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "manganato":
		data, err := s.service.ManganatoChapter(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "manganelo":
		data, err := s.service.ManganeloChapter(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangasee":
		data, err := s.service.MangaseeChapter(ctx, mangaId)
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

func (s *mangaControllerImpl) Detail(c *gin.Context) {

	ctx := c.Request.Context()
	source := c.Query("source")
	mangaId := c.Query("mangaid")

	switch source {
	case "mangabat":

		data, err := s.service.MangabatDetail(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangaread":

		data, err := s.service.MangareadDetail(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangatown":

		data, err := s.service.MangatownDetail(ctx, mangaId)
		if err != nil {
			fmt.Println(err)
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "maidmy":

		data, err := s.service.MaidmyDetail(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "asuracomic":

		data, err := s.service.AsuraComicDetail(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "manganato":
		data, err := s.service.ManganatoDetail(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "manganelo":
		data, err := s.service.ManganeloDetail(ctx, mangaId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangasee":
		data, err := s.service.MangaseeDetail(ctx, mangaId)
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

func (s *mangaControllerImpl) Search(c *gin.Context) {
	ctx := c.Request.Context()
	source := c.Query("source")
	title := strings.Replace(c.Query("title"), " ", "%20", -1)

	switch source {
	case "mangabat":

		data, err := s.service.MangabatSearch(ctx, title)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangaread":

		data, err := s.service.MangareadSearch(ctx, title)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangatown":

		data, err := s.service.MangatownSearch(ctx, title)
		if err != nil {
			fmt.Println(err)
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "maidmy":

		data, err := s.service.MaidmySearch(ctx, title)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "asuracomic":

		data, err := s.service.AsuraComicSearch(ctx, title)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "manganato":
		data, err := s.service.ManganatoSearch(ctx, title)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "manganelo":
		data, err := s.service.ManganeloSearch(ctx, title)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangasee":

		var err error
		if len(mangaseeJsonData) == 0 {
			mangaseeJsonData, err = s.service.MangaseeIndex(ctx)
			if err != nil {
				response.ErrorResponse(c, err, nil)
				return
			}
			mangaseeJsonDataExpired = time.Now()
		}

		if time.Now().UTC().After(mangaseeJsonDataExpired.Add(10 * time.Minute)) {
			log.Info().Msg("Expired Data, Refetching")
			mangaseeJsonData, err = s.service.MangaseeIndex(ctx)
			if err != nil {
				response.ErrorResponse(c, err, nil)
				return
			}
			mangaseeJsonDataExpired = time.Now()
		}

		data, err := utils.SearchIndex(mangaseeJsonData, title)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}

		if len(data) == 0 {
			response.ErrorResponse(c, utils.ERR_NOT_FOUND, utils.ERR_NOT_FOUND.Error())
			return
		}

		response.SuccesResponse(c, data)
		return

	default:
		response.ErrorResponse(c, utils.ERR_BAD_REQUEST, utils.ERR_BAD_REQUEST.Error())
		return
	}
}

func (s *mangaControllerImpl) Image(c *gin.Context) {

	ctx := c.Request.Context()
	source := c.Query("source")
	mangaId := c.Query("mangaid")
	chapterId := c.Query("chapterid")

	switch source {
	case "mangabat":

		data, err := s.service.MangabatImage(ctx, chapterId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangaread":

		data, err := s.service.MangareadImage(ctx, mangaId, chapterId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangatown":

		data, err := s.service.MangatownImage(ctx, mangaId, chapterId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "maidmy":

		data, err := s.service.MaidmyImage(ctx, chapterId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "asuracomic":

		data, err := s.service.AsuraComicImage(ctx, chapterId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "manganato":
		data, err := s.service.ManganatoImage(ctx, mangaId, chapterId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "manganelo":
		data, err := s.service.ManganeloImage(ctx, mangaId, chapterId)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}
		response.SuccesResponse(c, data)
		return

	case "mangasee":
		data, err := s.service.MangaseeImage(ctx, mangaId, chapterId)
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

func (s *mangaControllerImpl) MangabatProxy(c *gin.Context) {

	imageProxy := c.Query("id")
	req, err := http.NewRequest("GET", imageProxy, nil)
	if err != nil {
		response.ErrorResponse(c, err, nil)
		return
	}

	req.Header.Set("Referer", "https://m.mangabat.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		response.ErrorResponse(c, err, nil)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.ErrorResponse(c, err, nil)
		return
	}

	c.Writer.Write(body)

}

func (s *mangaControllerImpl) Redirect(c *gin.Context) {

	imageProxy := c.Query("id")
	sourceManga := c.Query("source")

	switch sourceManga {
	case "mangasee":

		var resp *http.Response
		var err error

		req, err := http.NewRequest("GET", imageProxy, nil)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}

		req.Header.Set("Referer", "https://mangasee123.com/")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {

			newImageProxy := strings.Replace(imageProxy, "official.lowee.us", "scans.lastation.us", -1)
			fmt.Println(newImageProxy)
			req, err := http.NewRequest("GET", newImageProxy, nil)
			if err != nil {
				response.ErrorResponse(c, err, nil)
				return
			}

			req.Header.Set("Referer", "https://mangasee123.com/")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")

			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				response.ErrorResponse(c, err, nil)
				return
			}

			defer resp.Body.Close()
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			response.ErrorResponse(c, err, nil)
			return
		}

		c.Writer.Write(body)

	default:
		response.ErrorResponse(c, utils.ERR_BAD_REQUEST, utils.ERR_BAD_REQUEST.Error())
		return
	}

}
