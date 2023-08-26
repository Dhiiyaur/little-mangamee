package service

import (
	"context"
	"little_mangamee/entity"
	log "little_mangamee/logger"
	"strings"

	"github.com/gocolly/colly"
)

type mangaServiceImpl struct{}

func NewMangaService() MangaService {
	return &mangaServiceImpl{}
}

func (m *mangaServiceImpl) MangabatIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error) {

	var returnData []entity.IndexData

	c := colly.NewCollector()
	c.OnHTML(".list-story-item", func(e *colly.HTMLElement) {

		tempLastChapter := strings.Split(e.ChildAttr("div > a:nth-child(2)", "href"), "-")
		tempMangaID := strings.Split(e.ChildAttr("a.item-img", "href"), "/")

		returnData = append(returnData, entity.IndexData{
			Title:          e.ChildAttr("img", "alt"),
			Id:             tempMangaID[len(tempMangaID)-1],
			Cover:          e.ChildAttr("a > img", "src"),
			LastChapter:    tempLastChapter[len(tempLastChapter)-1],
			OriginalServer: "https://m.mangabat.com/manga-list-all/" + pageNumber + "/",
		})
	})

	if err := c.Visit("https://m.mangabat.com/manga-list-all/" + pageNumber + "/"); err != nil {
		log.Info().Err(err)
		return nil, err
	}
	return returnData, nil
}

//

func (m *mangaServiceImpl) MangareadIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error) {

	var returnData []entity.IndexData

	c := colly.NewCollector()
	c.OnHTML(".page-item-detail.manga", func(e *colly.HTMLElement) {

		var coverImage string
		checkImage := strings.Split(e.ChildAttr("a > img", "data-src"), " ")
		if len(checkImage) < 2 {
			coverImage = e.ChildAttr("a > img", "data-src")
		} else {
			coverImage = checkImage[len(checkImage)-2]
		}
		returnData = append(returnData, entity.IndexData{

			Cover:          coverImage,
			Title:          e.ChildAttr("a", "title"),
			LastChapter:    strings.Split(e.ChildText("span.chapter.font-meta > a"), " ")[1],
			Id:             strings.Split(e.ChildAttr("a", "href"), "/")[4],
			OriginalServer: "https://www.mangaread.org/manga/?m_orderby=new-manga&page=" + pageNumber,
		})
	})

	if err := c.Visit("https://www.mangaread.org/manga/?m_orderby=new-manga&page=" + pageNumber); err != nil {
		log.Info().Err(err)
		return nil, err
	}

	return returnData, nil
}

//

func (m *mangaServiceImpl) MangatownIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error) {

	var returnData []entity.IndexData

	c := colly.NewCollector()
	c.OnHTML("body > section > article > div > div.manga_pic_content > ul.manga_pic_list > li", func(e *colly.HTMLElement) {
		var mangaId, lastChapter string

		mangaIdCheck := strings.Split(e.ChildAttr("a", "href"), "/")
		mangaId = mangaIdCheck[len(mangaIdCheck)-2]

		lastChapterCheck := strings.Split(e.ChildText("p.new_chapter"), " ")
		lastChapter = lastChapterCheck[len(lastChapterCheck)-1]

		mangaCoverCheck := strings.Replace(e.ChildAttr("a > img", "src"), "https://fmcdn.mangahere.com/", "http://fmcdn.mangatown.com/", -1)

		returnData = append(returnData, entity.IndexData{
			Id:             mangaId,
			Title:          e.ChildAttr("a", "title"),
			Cover:          mangaCoverCheck,
			LastChapter:    lastChapter,
			OriginalServer: "https://www.mangatown.com/hot/" + pageNumber + ".htm?wviews.za",
		})

	})
	if err := c.Visit("https://www.mangatown.com/hot/" + pageNumber + ".htm?wviews.za"); err != nil {
		log.Info().Err(err)
		return nil, err
	}

	return returnData, nil
}

//

func (m *mangaServiceImpl) MaidmyIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error) {

	var returnData []entity.IndexData

	c := colly.NewCollector()
	c.OnHTML("body > main > div > div.container > div.flexbox4 > div.flexbox4-item", func(e *colly.HTMLElement) {

		var checkLastChapter string

		tempLastChapter := strings.Split(e.ChildText("li > a"), " ")[1]
		if strings.Contains(tempLastChapter, "Ch.") {
			checkLastChapter = strings.Split(tempLastChapter, "C")[0]
		} else {
			checkLastChapter = tempLastChapter
		}
		returnData = append(returnData, entity.IndexData{
			Title:          e.ChildAttr("a", "title"),
			Id:             strings.Split(e.ChildAttr("a", "href"), "/")[4],
			Cover:          e.ChildAttr("img", "src"),
			LastChapter:    checkLastChapter,
			OriginalServer: "https://www.maid.my.id/page/" + pageNumber + "/",
		})
	})

	if err := c.Visit("https://www.maid.my.id/page/" + pageNumber + "/"); err != nil {
		log.Info().Err(err)
		return nil, err
	}

	return returnData, nil
}
