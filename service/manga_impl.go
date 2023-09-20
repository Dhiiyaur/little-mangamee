package service

import (
	"context"
	"encoding/json"
	"fmt"
	"little_mangamee/entity"
	log "little_mangamee/logger"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type mangaServiceImpl struct{}

func NewMangaService() MangaService {
	return &mangaServiceImpl{}
}

func (m *mangaServiceImpl) MangabatSearch(ctx context.Context, search string) ([]entity.SearchData, error) {

	var returnData []entity.SearchData

	c := colly.NewCollector()

	c.OnHTML(".list-story-item", func(e *colly.HTMLElement) {

		tempLastChapter := strings.Split(e.ChildAttr("div > a:nth-child(2)", "href"), "-")
		tempMangaID := strings.Split(e.ChildAttr("a.item-img", "href"), "/")

		returnData = append(returnData, entity.SearchData{
			Title:       e.ChildAttr("img", "alt"),
			Id:          tempMangaID[len(tempMangaID)-1],
			Cover:       e.ChildAttr("img", "src"),
			LastChapter: tempLastChapter[len(tempLastChapter)-1],
		})
	})

	if err := c.Visit("https://m.mangabat.com/search/manga/" + search); err != nil {
		return returnData, err
	}

	return returnData, nil
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
func (m *mangaServiceImpl) MangabatDetail(ctx context.Context, mangaId string) (entity.DetailData, error) {

	var returnData entity.DetailData
	var chapters []entity.Chapter

	c := colly.NewCollector()

	c.OnHTML(".panel-story-info", func(e *colly.HTMLElement) {

		returnData.Cover = e.ChildAttr("span > img", "src")
		returnData.Title = e.ChildText("div.story-info-right > h1")
		returnData.Summary = e.ChildText("div.panel-story-info-description")

	})

	c.OnHTML(".chapter-name", func(e *colly.HTMLElement) {

		tempMangaID := strings.Split(e.Attr("href"), "/")
		tempMangaName := strings.Split(e.Attr("href"), "-")

		chapters = append(chapters, entity.Chapter{
			Name: tempMangaName[len(tempMangaName)-1],
			Id:   tempMangaID[len(tempMangaID)-1],
		})

	})

	if err := c.Visit("https://readmangabat.com/" + mangaId + "/"); err != nil {
		return returnData, err
	}

	returnData.Chapters = chapters
	returnData.OriginalServer = "https://readmangabat.com/" + mangaId + "/"
	return returnData, nil
}
func (m *mangaServiceImpl) MangabatChapter(ctx context.Context, mangaId string) (entity.ChapterData, error) {

	var data entity.ChapterData
	var chapters []entity.Chapter

	c := colly.NewCollector()

	c.OnHTML(".chapter-name", func(e *colly.HTMLElement) {

		tempMangaID := strings.Split(e.Attr("href"), "/")
		tempMangaName := strings.Split(e.Attr("href"), "-")

		chapters = append(chapters, entity.Chapter{
			Name: tempMangaName[len(tempMangaName)-1],
			Id:   tempMangaID[len(tempMangaID)-1],
		})

	})

	if err := c.Visit("https://readmangabat.com/" + mangaId + "/"); err != nil {
		log.Info().Err(err)
		return data, err
	}

	data.Chapters = chapters
	data.OriginalServer = "https://readmangabat.com/" + mangaId
	return data, nil
}
func (m *mangaServiceImpl) MangabatImage(ctx context.Context, chapterId string) (entity.ImageData, error) {

	var returnData entity.ImageData
	var dataImages []entity.Image
	var name string

	c := colly.NewCollector()

	c.OnHTML(".img-content", func(e *colly.HTMLElement) {

		// fmt.Println(e.Attr("src"))
		api := "https://little-mangamee.vercel.app/api/manga/"
		dataImages = append(dataImages, entity.Image{
			// Image: fmt.Sprintf("%vproxy?id=%v", "https://api.mangamee.space/manga/", e.Attr("src")),
			Image: fmt.Sprintf("%vproxy?id=%v", api, e.Attr("src")),
		})

	})

	if err := c.Visit("https://readmangabat.com/" + chapterId + "/"); err != nil {
		return returnData, err
	}

	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	if strings.Contains(chapterId, "chap") {
		tmp := strings.Split(chapterId, "chap")
		name = re.FindAllString(tmp[len(tmp)-1], -1)[0]
	} else {
		name = re.FindAllString(chapterId, -1)[0]
	}

	returnData.OriginalServer = "https://readmangabat.com/" + chapterId + "/"
	returnData.ChapterName = name
	returnData.Images = dataImages

	return returnData, nil
}

//

func (m *mangaServiceImpl) MangareadSearch(ctx context.Context, search string) ([]entity.SearchData, error) {

	var returnData []entity.SearchData

	c := colly.NewCollector()
	c.OnHTML(".row.c-tabs-item__content", func(e *colly.HTMLElement) {

		var lastChapter string
		checkChapter := strings.Split(e.ChildText("span.font-meta.chapter > a"), " ")

		if len(checkChapter) > 2 {
			lastChapter = checkChapter[len(checkChapter)-2]
		} else {
			lastChapter = checkChapter[len(checkChapter)-1]
		}

		returnData = append(returnData, entity.SearchData{
			Cover:          e.ChildAttr("a > img", "data-src"),
			Title:          e.ChildAttr("a", "title"),
			LastChapter:    lastChapter,
			Id:             strings.Split(e.ChildAttr("a", "href"), "/")[4],
			OriginalServer: "https://www.mangaread.org/?s=" + search + "&post_type=wp-manga",
		})

	})

	if err := c.Visit("https://www.mangaread.org/?s=" + search + "&post_type=wp-manga"); err != nil {
		return returnData, err
	}

	return returnData, nil
}
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
func (m *mangaServiceImpl) MangareadDetail(ctx context.Context, mangaId string) (entity.DetailData, error) {

	var returnData entity.DetailData
	var chapters []entity.Chapter
	limit := 0

	c := colly.NewCollector()

	c.OnHTML(".post-title", func(e *colly.HTMLElement) {

		if limit == 0 {
			returnData.Title = strings.Split(e.ChildText("h1"), "  ")[0]
		}
		limit++
	})

	c.OnHTML(".summary_image", func(e *colly.HTMLElement) {
		returnData.Cover = e.ChildAttr("img", "data-src")
	})

	c.OnHTML(".summary__content", func(e *colly.HTMLElement) {
		returnData.Summary = e.ChildText("p")
	})

	c.OnHTML(".wp-manga-chapter", func(e *colly.HTMLElement) {

		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		tempName := strings.ReplaceAll(re.FindAllString(e.ChildText("a"), -1)[0], "-", "")

		chapters = append(chapters, entity.Chapter{
			Name: tempName,
			Id:   strings.Split(e.ChildAttr("a", "href"), "/")[5],
		})
	})

	if err := c.Visit("https://www.mangaread.org/manga/" + mangaId); err != nil {
		return returnData, err
	}

	returnData.Chapters = chapters
	returnData.OriginalServer = "https://www.mangaread.org/manga/" + mangaId

	return returnData, nil
}
func (m *mangaServiceImpl) MangareadChapter(ctx context.Context, mangaId string) (entity.ChapterData, error) {

	var data entity.ChapterData
	var chapters []entity.Chapter

	c := colly.NewCollector()

	c.OnHTML(".wp-manga-chapter", func(e *colly.HTMLElement) {

		re := regexp.MustCompile(`\d+`)
		tempName := re.FindAllString(e.ChildText("a"), -1)[0]

		chapters = append(chapters, entity.Chapter{
			Name: tempName,
			Id:   strings.Split(e.ChildAttr("a", "href"), "/")[5],
		})
	})

	if err := c.Visit("https://www.mangaread.org/manga/" + mangaId); err != nil {
		log.Info().Err(err)
		return data, err
	}

	data.Chapters = chapters
	data.OriginalServer = "https://www.mangaread.org/manga/" + mangaId
	return data, nil
}
func (m *mangaServiceImpl) MangareadImage(ctx context.Context, mangaId string, chapterId string) (entity.ImageData, error) {

	var returnData entity.ImageData
	var dataImages []entity.Image
	c := colly.NewCollector()

	c.OnHTML(".wp-manga-chapter-img", func(e *colly.HTMLElement) {

		dataImages = append(dataImages, entity.Image{
			Image: "https://" + strings.Split(e.Attr("data-src"), "//")[1],
		})

	})
	if err := c.Visit("https://www.mangaread.org/manga/" + mangaId + "/" + chapterId); err != nil {
		return returnData, err
	}

	re := regexp.MustCompile(`\d+`)

	returnData.OriginalServer = "https://www.mangaread.org/manga/" + mangaId + "/" + chapterId
	returnData.ChapterName = re.FindAllString(chapterId, -1)[0]
	returnData.Images = dataImages
	return returnData, nil
}

//

func (m *mangaServiceImpl) MangatownSearch(ctx context.Context, search string) ([]entity.SearchData, error) {

	var returnData []entity.SearchData
	c := colly.NewCollector()
	c.OnHTML(".manga_pic_list > li", func(e *colly.HTMLElement) {

		mangaCoverCheck := strings.Replace(e.ChildAttr("a.manga_cover > img", "src"), "https://fmcdn.mangahere.com/", "http://fmcdn.mangatown.com/", -1)
		returnData = append(returnData, entity.SearchData{
			Cover:          mangaCoverCheck,
			Title:          e.ChildAttr("a.manga_cover", "title"),
			Id:             strings.Split(e.ChildAttr("a.manga_cover", "href"), "/")[2],
			OriginalServer: "https://www.mangatown.com/search?name=" + search,
		})

	})

	if err := c.Visit("https://www.mangatown.com/search?name=" + search); err != nil {
		return returnData, err
	}

	return returnData, nil
}
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
func (m *mangaServiceImpl) MangatownDetail(ctx context.Context, mangaId string) (entity.DetailData, error) {

	var returnData entity.DetailData
	var chapters []entity.Chapter

	c := colly.NewCollector()

	c.OnHTML(".article_content", func(e *colly.HTMLElement) {

		mangaCoverCheck := strings.Replace(e.ChildAttr("div.detail_info.clearfix > img", "src"), "https://fmcdn.mangahere.com/", "http://fmcdn.mangatown.com/", -1)

		returnData.Title = e.ChildText("h1.title-top")
		returnData.Cover = mangaCoverCheck
		returnData.Summary = e.ChildText("div.detail_info.clearfix > ul > li > span")

	})

	c.OnHTML(".chapter_list > li", func(e *colly.HTMLElement) {

		var chapterName string

		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		arr := re.FindAllString(e.ChildText("a"), -1)
		if len(arr) != 0 {
			chapterName = arr[len(arr)-1]
		} else {
			chapterName = "0"
		}

		chapters = append(chapters, entity.Chapter{
			Id:   strings.Split(e.ChildAttr("a", "href"), "/")[3],
			Name: chapterName,
		})
	})

	if err := c.Visit("https://www.mangatown.com/manga/" + mangaId); err != nil {
		return returnData, err
	}

	returnData.Chapters = chapters
	returnData.OriginalServer = "https://www.mangatown.com/manga/" + mangaId

	return returnData, nil
}
func (m *mangaServiceImpl) MangatownChapter(ctx context.Context, mangaId string) (entity.ChapterData, error) {

	var data entity.ChapterData
	var chapters []entity.Chapter

	c := colly.NewCollector()

	c.OnHTML(".chapter_list > li", func(e *colly.HTMLElement) {

		var chapterName string

		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		arr := re.FindAllString(e.ChildText("a"), -1)
		if len(arr) != 0 {
			chapterName = arr[len(arr)-1]
		} else {
			chapterName = "0"
		}

		chapters = append(chapters, entity.Chapter{
			Id:   strings.Split(e.ChildAttr("a", "href"), "/")[3],
			Name: chapterName,
		})
	})

	if err := c.Visit("https://www.mangatown.com/manga/" + mangaId); err != nil {
		log.Info().Err(err)
		return data, err
	}

	data.Chapters = chapters
	data.OriginalServer = "https://www.mangatown.com/manga/" + mangaId
	return data, nil
}
func (m *mangaServiceImpl) MangatownImage(ctx context.Context, mangaId string, chapterId string) (entity.ImageData, error) {

	var returnData entity.ImageData
	var dataImages []entity.Image
	var link string

	c := colly.NewCollector()

	c.OnHTML(".read_img", func(e *colly.HTMLElement) {

		mangaCoverCheck := strings.Replace(e.ChildAttr("img", "src"), "zjcdn.mangahere.org", "fmcdn.mangatown.com", -1)
		link = "https:" + mangaCoverCheck
	})

	if err := c.Visit("https://www.mangatown.com/manga/" + mangaId + "/" + chapterId + "/"); err != nil {
		return returnData, err
	}

	baseLink, imageLink := returnLastSliceAndJoinLink(link)
	imageExtension, frontRawData, loopData := getRawImageData(imageLink)

	for i := 0; i < 100; i++ {
		tempNumber := loopData + i
		if tempNumber < 10 {
			a := fmt.Sprintf("%v00%v.%v", frontRawData, strconv.Itoa(tempNumber), imageExtension)
			joinImageLink := fmt.Sprintf("%v/%v", baseLink, a)
			dataImages = append(dataImages, entity.Image{
				Image: joinImageLink,
			})

		} else if tempNumber < 100 && tempNumber > 9 {
			a := fmt.Sprintf("%v0%v.%v", frontRawData, strconv.Itoa(tempNumber), imageExtension)
			joinImageLink := fmt.Sprintf("%v/%v", baseLink, a)
			dataImages = append(dataImages, entity.Image{
				Image: joinImageLink,
			})

		} else if tempNumber < 1000 && tempNumber > 99 {
			a := fmt.Sprintf("%v%v.%v", frontRawData, strconv.Itoa(tempNumber), imageExtension)
			joinImageLink := fmt.Sprintf("%v/%v", baseLink, a)
			dataImages = append(dataImages, entity.Image{
				Image: joinImageLink,
			})

		} else if tempNumber < 10000 && tempNumber > 999 {
			a := fmt.Sprintf("%v%v.%v", frontRawData, strconv.Itoa(tempNumber), imageExtension)
			joinImageLink := fmt.Sprintf("%v/%v", baseLink, a)
			dataImages = append(dataImages, entity.Image{
				Image: joinImageLink,
			})

		}
	}

	re := regexp.MustCompile(`\d+`)
	returnData.OriginalServer = "https://www.mangatown.com/manga/" + mangaId + "/" + chapterId + "/"
	returnData.ChapterName = re.FindAllString(chapterId, -1)[0]
	returnData.Images = dataImages
	return returnData, nil
}
func returnLastSliceAndJoinLink(s string) (string, string) {
	slice := strings.Split(s, "/")
	return strings.Join(slice[:len(slice)-1], "/"), slice[len(slice)-1]
}
func getRawImageData(s string) (string, string, int) {

	var imageExtension, frontRawData string
	var loopData int

	a := strings.Split(s, ".")
	imageExtension = a[len(a)-1]

	if strings.Contains(s, "_") {
		b := strings.Split(a[0], "_")
		loopData, _ = strconv.Atoi(b[len(b)-1])

		if len(b) > 2 {
			frontRawData = fmt.Sprintf("%v_%v_", b[0], b[1])
		} else {
			frontRawData = fmt.Sprintf("%v_", b[0])
		}

	} else {
		frontRawData = a[0][0:1]
		loopData, _ = strconv.Atoi(a[0][1:])
	}

	return imageExtension, frontRawData, loopData
}

func (m *mangaServiceImpl) MaidmySearch(ctx context.Context, search string) ([]entity.SearchData, error) {

	var returnData []entity.SearchData

	c := colly.NewCollector()

	c.OnHTML("body > main > div > div > div.flexbox2 > div.flexbox2-item", func(e *colly.HTMLElement) {

		var checkLastChapter string

		tempLastChapter := strings.Split(e.ChildText("div.season"), " ")

		if len(tempLastChapter) > 1 {
			checkLastChapter = tempLastChapter[1]
		}

		returnData = append(returnData, entity.SearchData{
			Title:       e.ChildAttr("a", "title"),
			Id:          strings.Split(e.ChildAttr("a", "href"), "/")[4],
			Cover:       e.ChildAttr("img", "src"),
			LastChapter: checkLastChapter,
		})
	})

	if err := c.Visit("https://www.maid.my.id/?s=" + search); err != nil {
		return returnData, err
	}

	return returnData, nil
}
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
func (m *mangaServiceImpl) MaidmyDetail(ctx context.Context, mangaId string) (entity.DetailData, error) {

	var returnData entity.DetailData
	var chapters []entity.Chapter

	c := colly.NewCollector()

	c.OnHTML(".series-thumb", func(e *colly.HTMLElement) {
		returnData.Cover = e.ChildAttr(`img`, "src")
	})

	c.OnHTML(".series-title", func(e *colly.HTMLElement) {
		returnData.Title = e.ChildText(`h2`)
	})

	c.OnHTML(".series-synops", func(e *colly.HTMLElement) {
		returnData.Summary = e.Text

	})

	c.OnHTML(".flexch-infoz", func(e *colly.HTMLElement) {

		var chapterName string
		tempChapterName := e.ChildAttr(`a`, "title")

		if strings.Contains(tempChapterName, "Bahasa Indonesia") {
			a := strings.Split(tempChapterName, "Bahasa Indonesia")
			b := strings.Split(a[len(a)-2], " ")
			chapterName = fmt.Sprintf("%v %v", b[len(b)-3], b[len(b)-2])

		} else {
			a := strings.Split(tempChapterName, " ")
			chapterName = fmt.Sprintf("%v %v", a[len(a)-2], a[len(a)-1])
		}

		chapters = append(chapters, entity.Chapter{
			Name: chapterName,
			Id:   strings.Split(e.ChildAttr(`a`, "href"), "/")[3],
		})

	})

	if err := c.Visit("https://www.maid.my.id/manga/" + mangaId + "/"); err != nil {
		return returnData, err
	}

	returnData.Chapters = chapters
	returnData.OriginalServer = "https://www.maid.my.id/manga/" + mangaId + "/"
	return returnData, nil
}
func (m *mangaServiceImpl) MaidmyChapter(ctx context.Context, mangaId string) (entity.ChapterData, error) {

	var data entity.ChapterData
	var chapters []entity.Chapter

	c := colly.NewCollector()

	c.OnHTML(".flexch-infoz", func(e *colly.HTMLElement) {

		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		tempName := strings.ReplaceAll(re.FindAllString(e.ChildAttr(`a`, "title"), -1)[0], "-", "")
		chapters = append(chapters, entity.Chapter{
			Name: tempName,
			Id:   strings.Split(e.ChildAttr(`a`, "href"), "/")[3],
		})

	})

	if err := c.Visit("https://www.maid.my.id/manga/" + mangaId + "/"); err != nil {
		log.Info().Err(err)
		return data, err
	}

	data.Chapters = chapters
	data.OriginalServer = "https://www.maid.my.id/manga/" + mangaId
	return data, nil
}
func (m *mangaServiceImpl) MaidmyImage(ctx context.Context, chapterId string) (entity.ImageData, error) {

	var returnData entity.ImageData
	var dataImages []entity.Image
	c := colly.NewCollector()

	c.OnHTML(".reader-area img", func(e *colly.HTMLElement) {
		dataImages = append(dataImages, entity.Image{
			Image: e.Attr("src"),
		})

	})

	if err := c.Visit("https://www.maid.my.id/" + chapterId); err != nil {
		log.Info().Err(err)
		return returnData, err
	}

	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	returnData.Images = dataImages
	returnData.ChapterName = re.FindAllString(chapterId, -1)[0]
	returnData.OriginalServer = "https://www.maid.my.id/" + chapterId
	return returnData, nil
}

func (m *mangaServiceImpl) AsuraComicIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error) {

	var returnData []entity.IndexData

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("#content > div.wrapper > div.postbody > div.bixbox > div.mrgn > div.listupd > div.bs > div.bsx", func(e *colly.HTMLElement) {
		returnData = append(returnData, entity.IndexData{
			Title:          e.ChildText("a > div.bigor > div.tt"),
			Id:             strings.Split(e.ChildAttr("a", "href"), "/")[4],
			Cover:          e.ChildAttr("a > div.limit > img", "src"),
			LastChapter:    e.ChildText("a > div.bigor > div.adds > div.epxs"),
			OriginalServer: fmt.Sprintf("https://asuracomics.com/manga/?page=%v&order=update", pageNumber),
		})
	})

	if err := c.Visit(fmt.Sprintf("https://asuracomics.com/manga/?page=%v&order=update", pageNumber)); err != nil {
		log.Info().Err(err)
		return nil, err
	}

	return returnData, nil
}

func (m *mangaServiceImpl) AsuraComicSearch(ctx context.Context, search string) ([]entity.SearchData, error) {

	var returnData []entity.SearchData

	pageCount := 3
	currentPage := 1

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)
	c.OnHTML("#content > div.wrapper > div.postbody > div.bixbox > div.listupd > div.bs > div.bsx", func(e *colly.HTMLElement) {
		returnData = append(returnData, entity.SearchData{
			Title:          e.ChildText("a > div.bigor > div.tt"),
			Id:             strings.Split(e.ChildAttr("a", "href"), "/")[4],
			Cover:          e.ChildAttr("a > div.limit > img", "src"),
			LastChapter:    e.ChildText("a > div.bigor > div.adds > div.epxs"),
			OriginalServer: fmt.Sprintf("https://asuracomics.com/page/%v/?s=solo", currentPage),
		})
	})

	for currentPage <= pageCount {
		if err := c.Visit(fmt.Sprintf("https://asuracomics.com/page/%v/?s=%v", currentPage, search)); err != nil {
			log.Info().Err(err)
			return returnData, nil
		}
		currentPage++
	}

	return returnData, nil

}

func (m *mangaServiceImpl) AsuraComicDetail(ctx context.Context, mangaId string) (entity.DetailData, error) {

	var returnData entity.DetailData
	var chapters []entity.Chapter

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	c.OnHTML("div.bixbox.animefull", func(e *colly.HTMLElement) {

		returnData.Title = e.ChildText("div.bigcontent > div.infox > h1.entry-title")
		returnData.Cover = e.ChildAttr("div.bigcontent > div.thumbook > div.thumb > img", "src")
		returnData.OriginalServer = fmt.Sprintf("https://asuracomics.com/manga/%v", mangaId)
		returnData.Summary = e.ChildText("div.bigcontent > div.infox > div.wd-full > div.entry-content.entry-content-single")

	})

	c.OnHTML("#chapterlist > ul > li > div.chbox", func(e *colly.HTMLElement) {

		fmt.Println("link", strings.Split(e.ChildAttr("div.eph-num > a", "href"), "/")[3])
		fmt.Println("chapname", re.FindAllString(e.ChildText("div.eph-num > a > span.chapternum"), -1)[0])

		chapters = append(chapters, entity.Chapter{
			Id:   strings.Split(e.ChildAttr("div.eph-num > a", "href"), "/")[3],
			Name: re.FindAllString(e.ChildText("div.eph-num > a > span.chapternum"), -1)[0],
		})
	})

	if err := c.Visit(fmt.Sprintf("https://asuracomics.com/manga/%v", mangaId)); err != nil {
		log.Info().Err(err)
		return returnData, nil
	}

	returnData.Chapters = chapters
	return returnData, nil
}

func (m *mangaServiceImpl) AsuraComicChapter(ctx context.Context, mangaId string) (entity.ChapterData, error) {

	var data entity.ChapterData
	var chapters []entity.Chapter

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	c.OnHTML("#chapterlist > ul > li > div.chbox", func(e *colly.HTMLElement) {

		fmt.Println("link", strings.Split(e.ChildAttr("div.eph-num > a", "href"), "/")[3])
		fmt.Println("chapname", re.FindAllString(e.ChildText("div.eph-num > a > span.chapternum"), -1)[0])

		chapters = append(chapters, entity.Chapter{
			Id:   strings.Split(e.ChildAttr("div.eph-num > a", "href"), "/")[3],
			Name: re.FindAllString(e.ChildText("div.eph-num > a > span.chapternum"), -1)[0],
		})
	})

	if err := c.Visit(fmt.Sprintf("https://asuracomics.com/manga/%v", mangaId)); err != nil {
		log.Info().Err(err)
		return data, nil
	}

	data.OriginalServer = fmt.Sprintf("https://asuracomics.com/manga/%v", mangaId)
	data.Chapters = chapters
	return data, nil
}

func (m *mangaServiceImpl) AsuraComicImage(ctx context.Context, chapterId string) (entity.ImageData, error) {

	var returnData entity.ImageData
	var dataImages []entity.Image

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)
	c.OnHTML("#readerarea > p > img", func(e *colly.HTMLElement) {
		dataImages = append(dataImages, entity.Image{
			Image: e.Attr("src"),
		})
	})

	if err := c.Visit(fmt.Sprintf("https://asuracomics.com/%v", chapterId)); err != nil {
		log.Info().Err(err)
		return returnData, nil
	}

	returnData.Images = dataImages
	returnData.OriginalServer = fmt.Sprintf("https://asuracomics.com/%v", chapterId)
	returnData.ChapterName = chapterId

	return returnData, nil
}

func (m *mangaServiceImpl) ManganatoSearch(ctx context.Context, search string) ([]entity.SearchData, error) {

	var returnData []entity.SearchData

	pageCount := 3
	currentPage := 1

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-search-story > div.search-story-item", func(e *colly.HTMLElement) {

		returnData = append(returnData, entity.SearchData{
			Title:          e.ChildText("div.item-right > h3 > a"),
			Id:             strings.Split(e.ChildAttr("a", "href"), "/")[3],
			Cover:          e.ChildAttr("a > img", "src"),
			LastChapter:    e.ChildText("div.item-right > a:nth-child(2)"),
			OriginalServer: fmt.Sprintf("https://manganato.com/search/story/%v?page=%v", search, currentPage),
		})
	})

	for currentPage <= pageCount {
		if err := c.Visit(fmt.Sprintf("https://manganato.com/search/story/%v?page=%v", search, currentPage)); err != nil {
			log.Info().Err(err)
			return returnData, nil
		}
		currentPage++
	}

	return returnData, nil

}

func (m *mangaServiceImpl) ManganatoIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error) {

	var returnData []entity.IndexData

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container.container-main > div.panel-content-genres > div.content-genres-item", func(e *colly.HTMLElement) {
		returnData = append(returnData, entity.IndexData{
			Title:          e.ChildText("div.genres-item-info > h3 > a"),
			Id:             strings.Split(e.ChildAttr("a", "href"), "/")[3],
			Cover:          e.ChildAttr("a > img", "src"),
			LastChapter:    e.ChildText("div.genres-item-info > a.genres-item-chap.text-nowrap.a-h"),
			OriginalServer: fmt.Sprintf("https://manganato.com/advanced_search?s=all&g_e=_41_&page=%v", pageNumber),
		})
	})

	if err := c.Visit(fmt.Sprintf("https://manganato.com/advanced_search?s=all&g_e=_41_&page=%v", pageNumber)); err != nil {
		log.Info().Err(err)
		return nil, err
	}

	return returnData, nil
}

func (m *mangaServiceImpl) ManganatoDetail(ctx context.Context, mangaId string) (entity.DetailData, error) {

	var returnData entity.DetailData
	var chapters []entity.Chapter

	// re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-info", func(e *colly.HTMLElement) {

		tmpSummary := strings.Split(e.ChildText("div.panel-story-info-description"), ".-")
		removeText := strings.Split(tmpSummary[0], "Description :")
		if len(removeText) > 1 {
			returnData.Summary = removeText[1]
		} else {
			returnData.Summary = removeText[0]
		}

		returnData.Cover = e.ChildAttr("span.info-image > img", "src")
		returnData.Title = e.ChildText("div.story-info-right > h1")
	})

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-chapter-list > ul > li.a-h", func(e *colly.HTMLElement) {
		chapters = append(chapters, entity.Chapter{
			Name: e.ChildText("a"),
			Id:   strings.Split(e.ChildAttr("a", "href"), "/")[4],
		})

	})

	if err := c.Visit(fmt.Sprintf("https://chapmanganato.com/%v", mangaId)); err != nil {
		log.Info().Err(err)
		return returnData, err
	}

	returnData.Chapters = chapters
	returnData.OriginalServer = fmt.Sprintf("https://chapmanganato.com/%v", mangaId)
	return returnData, nil

}

func (m *mangaServiceImpl) ManganatoChapter(ctx context.Context, mangaId string) (entity.ChapterData, error) {

	var data entity.ChapterData
	var chapters []entity.Chapter
	// re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-chapter-list > ul > li.a-h", func(e *colly.HTMLElement) {
		chapters = append(chapters, entity.Chapter{
			// Name: re.FindAllString(e.ChildText("a"), -1)[0],
			Name: e.ChildText("a"),
			Id:   strings.Split(e.ChildAttr("a", "href"), "/")[4],
		})

	})

	if err := c.Visit(fmt.Sprintf("https://chapmanganato.com/%v", mangaId)); err != nil {
		log.Info().Err(err)
		return data, err
	}

	data.Chapters = chapters
	data.OriginalServer = fmt.Sprintf("https://chapmanganato.com/%v", mangaId)
	return data, nil
}

func (m *mangaServiceImpl) ManganatoImage(ctx context.Context, mangaId string, chapterId string) (entity.ImageData, error) {

	var returnData entity.ImageData
	var dataImages []entity.Image

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container-chapter-reader > img", func(e *colly.HTMLElement) {
		dataImages = append(dataImages, entity.Image{
			Image: e.Attr("src"),
		})

	})

	if err := c.Visit(fmt.Sprintf("https://chapmanganato.com/%v/%v", mangaId, chapterId)); err != nil {
		log.Info().Err(err)
		return returnData, err
	}

	returnData.Images = dataImages
	returnData.OriginalServer = fmt.Sprintf("https://chapmanganato.com/%v/%v", mangaId, chapterId)
	returnData.ChapterName = chapterId

	return returnData, nil

}

func (m *mangaServiceImpl) ManganeloIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error) {

	var returnData []entity.IndexData
	// re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container.container-main > div.panel-content-genres > div.content-genres-item", func(e *colly.HTMLElement) {

		// tmpLastChapter := re.FindAllString(e.ChildText("div.genres-item-info > a.genres-item-chap.text-nowrap.a-h"), -1)
		// var lastChapter string
		// if len(tmpLastChapter) > 0 {
		// 	lastChapter = tmpLastChapter[1]
		// } else {
		// 	lastChapter = tmpLastChapter[0]
		// }

		lastChapter := e.ChildText("div.genres-item-info > a.genres-item-chap.text-nowrap.a-h")
		returnData = append(returnData, entity.IndexData{
			Title:          e.ChildText("div.genres-item-info > h3"),
			Id:             strings.Split(e.ChildAttr("a", "href"), "/")[3],
			Cover:          e.ChildAttr("a > img", "src"),
			LastChapter:    lastChapter,
			OriginalServer: fmt.Sprintf("https://m.manganelo.com/advanced_search?s=all&g_e=_41_&page=%v", pageNumber),
		})
	})

	if err := c.Visit(fmt.Sprintf("https://m.manganelo.com/advanced_search?s=all&g_e=_41_&page=%v", pageNumber)); err != nil {
		log.Info().Err(err)
		return nil, err
	}

	return returnData, nil
}

func (m *mangaServiceImpl) ManganeloSearch(ctx context.Context, search string) ([]entity.SearchData, error) {

	var returnData []entity.SearchData

	pageCount := 3
	currentPage := 1

	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-search-story > div.search-story-item", func(e *colly.HTMLElement) {

		returnData = append(returnData, entity.SearchData{
			Title:          e.ChildText("div.item-right > h3 > a"),
			Id:             strings.Split(e.ChildAttr("a", "href"), "/")[3],
			Cover:          e.ChildAttr("a > img", "src"),
			LastChapter:    re.FindAllString(e.ChildText("div.item-right > a:nth-child(2)"), -1)[0],
			OriginalServer: fmt.Sprintf("https://m.manganelo.com/search/story/%v?page=%v", search, currentPage),
		})
	})

	for currentPage <= pageCount {
		if err := c.Visit(fmt.Sprintf("https://m.manganelo.com/search/story/%v?page=%v", search, currentPage)); err != nil {
			log.Info().Err(err)
			return returnData, nil
		}
		currentPage++
	}

	return returnData, nil

}

func (m *mangaServiceImpl) ManganeloDetail(ctx context.Context, mangaId string) (entity.DetailData, error) {

	var returnData entity.DetailData
	var chapters []entity.Chapter

	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-info", func(e *colly.HTMLElement) {

		removeText := strings.Split(e.ChildText("div.panel-story-info-description"), "Description :")
		if len(removeText) > 1 {
			returnData.Summary = removeText[1]
		} else {
			returnData.Summary = removeText[0]
		}

		returnData.Cover = e.ChildAttr("span.info-image > img.img-loading", "src")
		returnData.Title = e.ChildText("div.story-info-right > h1")
	})

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-chapter-list > ul > li.a-h", func(e *colly.HTMLElement) {
		chapters = append(chapters, entity.Chapter{
			Name: re.FindAllString(e.ChildText("a"), -1)[0],
			Id:   strings.Split(e.ChildAttr("a", "href"), "/")[4],
		})
	})

	if err := c.Visit(fmt.Sprintf("https://m.manganelo.com/%v", mangaId)); err != nil {
		log.Info().Err(err)
		return returnData, err
	}

	returnData.OriginalServer = fmt.Sprintf("https://m.manganelo.com/%v", mangaId)
	if len(chapters) == 0 {

		c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-info", func(e *colly.HTMLElement) {

			removeText := strings.Split(e.ChildText("div.panel-story-info-description"), "Description :")
			if len(removeText) > 1 {
				returnData.Summary = removeText[1]
			} else {
				returnData.Summary = removeText[0]
			}

			returnData.Cover = e.ChildAttr("span.info-image > img.img-loading", "src")
			returnData.Title = e.ChildText("div.story-info-right > h1")
		})

		c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-chapter-list > ul > li.a-h", func(e *colly.HTMLElement) {
			chapters = append(chapters, entity.Chapter{
				Name: re.FindAllString(e.ChildText("a"), -1)[0],
				Id:   strings.Split(e.ChildAttr("a", "href"), "/")[4],
			})
		})

		if err := c.Visit(fmt.Sprintf("https://chapmanganelo.com/%v", mangaId)); err != nil {
			log.Info().Err(err)
			return returnData, err
		}
		returnData.OriginalServer = fmt.Sprintf("https://chapmanganelo.com/%v", mangaId)
	}

	returnData.Chapters = chapters
	return returnData, nil

}

func (m *mangaServiceImpl) ManganeloChapter(ctx context.Context, mangaId string) (entity.ChapterData, error) {

	var data entity.ChapterData
	var chapters []entity.Chapter
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-chapter-list > ul > li.a-h", func(e *colly.HTMLElement) {
		chapters = append(chapters, entity.Chapter{
			Name: re.FindAllString(e.ChildText("a"), -1)[0],
			Id:   strings.Split(e.ChildAttr("a", "href"), "/")[4],
		})
	})

	if err := c.Visit(fmt.Sprintf("https://manganelo.com/%v", mangaId)); err != nil {
		log.Info().Err(err)
		return data, err
	}
	data.OriginalServer = fmt.Sprintf("https://manganelo.com/%v", mangaId)

	if len(chapters) == 0 {
		c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-chapter-list > ul > li.a-h", func(e *colly.HTMLElement) {
			chapters = append(chapters, entity.Chapter{
				Name: re.FindAllString(e.ChildText("a"), -1)[0],
				Id:   strings.Split(e.ChildAttr("a", "href"), "/")[4],
			})
		})

		if err := c.Visit(fmt.Sprintf("https://chapmanganelo.com/%v", mangaId)); err != nil {
			log.Info().Err(err)
			return data, err
		}
		data.OriginalServer = fmt.Sprintf("https://chapmanganelo.com/%v", mangaId)

	}

	data.Chapters = chapters
	return data, nil
}

func (m *mangaServiceImpl) ManganeloImage(ctx context.Context, mangaId string, chapterId string) (entity.ImageData, error) {

	var returnData entity.ImageData
	var dataImages []entity.Image

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container-chapter-reader > img.reader-content", func(e *colly.HTMLElement) {
		dataImages = append(dataImages, entity.Image{
			Image: e.Attr("src"),
		})

	})

	if err := c.Visit(fmt.Sprintf("https://chapmanganelo.com/%v/%v", mangaId, chapterId)); err != nil {
		log.Info().Err(err)
		return returnData, err
	}

	returnData.Images = dataImages
	returnData.OriginalServer = fmt.Sprintf("https://chapmanganelo.com/%v/%v", mangaId, chapterId)
	returnData.ChapterName = chapterId

	return returnData, nil

}

type customDate struct {
	time.Time
}

func (ct *customDate) UnmarshalJSON(b []byte) (err error) {

	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}

	if s == "0" {
		ct.Time = time.Time{}
		return
	}

	ct.Time, err = time.Parse(time.RFC3339, s)
	return
}

func (m *mangaServiceImpl) MangaseeIndex(ctx context.Context) ([]entity.IndexData, error) {

	var match string

	type MangaseeData struct {
		Id      string     `json:"i"`
		Title   string     `json:"s"`
		Date    customDate `json:"ls"`
		Chapter string     `json:"l"`
	}

	c := colly.NewCollector()
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text // Get the text content of the <script> tag
		re := regexp.MustCompile(`vm.Directory = \[(.*?)\];`)
		match = re.FindString(scriptContent)
	})

	if err := c.Visit("https://www.mangasee123.com/search/?sort=lt&desc=true"); err != nil {
		log.Info().Err(err)
		return nil, err
	}

	jsonData := strings.Split(match, "vm.Directory = ")
	jsonData = strings.Split(jsonData[1], "];")
	cleanData := jsonData[0] + "]"

	var mangaData []MangaseeData
	err := json.Unmarshal([]byte(cleanData), &mangaData)
	if err != nil {
		log.Info().Err(err)
		return nil, err
	}

	sort.Slice(mangaData, func(i, j int) bool { return mangaData[i].Date.Time.After(mangaData[j].Date.Time) })
	var returnData []entity.IndexData

	for _, v := range mangaData {

		var lastChapter string
		if v.Chapter == "N/A" {
			lastChapter = "-"
		} else {
			lastChapter = strings.TrimPrefix(v.Chapter, "10")
			lastChapter = strings.TrimSuffix(lastChapter, "0")
		}

		returnData = append(returnData, entity.IndexData{
			Title: v.Title,
			// Id:    fmt.Sprintf("https://www.mangasee123.com/manga/%v", v.Id),
			Id:             v.Id,
			Cover:          fmt.Sprintf("https://temp.compsci88.com/cover/%v", v.Id),
			LastChapter:    lastChapter,
			OriginalServer: "https://www.mangasee123.com/search/?sort=lt&desc=true",
		})
	}

	return returnData, nil
}

func (m *mangaServiceImpl) MangaseeDetail(ctx context.Context, mangaId string) (entity.DetailData, error) {

	var returnData entity.DetailData
	var chapters []entity.Chapter
	var cover, title, summary, baseChapter string

	type MangaseeDataChapter struct {
		Chapter string `json:"Chapter"`
	}

	c := colly.NewCollector()
	c.OnHTML("body > div.container.MainContainer > div > div > div > div > div.row", func(e *colly.HTMLElement) {
		if cover == "" {
			cover = e.ChildAttr("div.col-md-3.col-sm-4.col-3.top-5 > img.img-fluid.bottom-5", "src")
		}

		if title == "" {
			title = e.ChildText("div.bottom-10")
		}

		summary = summary + e.ChildText("div.top-5.Content")

	})

	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text // Get the text content of the <script> tag
		re := regexp.MustCompile(`vm.Chapters = \[(.*?)\];`)
		baseChapter = re.FindString(scriptContent)
	})

	if err := c.Visit(fmt.Sprintf("https://www.mangasee123.com/manga/%v", mangaId)); err != nil {
		log.Info().Err(err)
		return returnData, err
	}

	baseChapter = strings.Split(baseChapter, "vm.Chapters = ")[1]
	baseChapter = strings.Split(baseChapter, ";")[0]
	var mangaData []MangaseeDataChapter
	err := json.Unmarshal([]byte(baseChapter), &mangaData)
	if err != nil {
		log.Info().Err(err)
		return returnData, err
	}

	for _, v := range mangaData {

		v.Chapter = strings.TrimPrefix(v.Chapter, "10")
		v.Chapter = strings.TrimSuffix(v.Chapter, "0")

		chapters = append(chapters, entity.Chapter{
			Name: v.Chapter,
			Id:   v.Chapter,
		})
	}

	returnData.Chapters = chapters
	returnData.Cover = cover
	returnData.OriginalServer = fmt.Sprintf("https://www.mangasee123.com/manga/%v", mangaId)
	returnData.Summary = summary
	returnData.Title = title

	return returnData, nil
}

func (m *mangaServiceImpl) MangaseeChapter(ctx context.Context, mangaId string) (entity.ChapterData, error) {
	var data entity.ChapterData
	var chapters []entity.Chapter

	var baseChapter string

	type MangaseeDataChapter struct {
		Chapter string `json:"Chapter"`
	}

	c := colly.NewCollector()
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text // Get the text content of the <script> tag
		re := regexp.MustCompile(`vm.Chapters = \[(.*?)\];`)
		baseChapter = re.FindString(scriptContent)
	})

	if err := c.Visit(fmt.Sprintf("https://www.mangasee123.com/manga/%v", mangaId)); err != nil {
		log.Info().Err(err)
		return data, err
	}

	baseChapter = strings.Split(baseChapter, "vm.Chapters = ")[1]
	baseChapter = strings.Split(baseChapter, ";")[0]
	var mangaData []MangaseeDataChapter
	err := json.Unmarshal([]byte(baseChapter), &mangaData)
	if err != nil {
		log.Info().Err(err)
		return data, err
	}

	for _, v := range mangaData {

		v.Chapter = strings.TrimPrefix(v.Chapter, "10")
		v.Chapter = strings.TrimSuffix(v.Chapter, "0")

		chapters = append(chapters, entity.Chapter{
			Name: v.Chapter,
			Id:   v.Chapter,
		})
	}

	data.Chapters = chapters
	data.OriginalServer = fmt.Sprintf("https://www.mangasee123.com/manga/%v", mangaId)
	return data, nil
}

func generateNumber(number int) string {
	digit := "%03d"
	return fmt.Sprintf(digit, number)
}

func (m *mangaServiceImpl) MangaseeImage(ctx context.Context, mangaId string, chapterId string) (entity.ImageData, error) {

	var returnData entity.ImageData
	var dataImages []entity.Image

	var baseImages string

	type MangaseeDataImage struct {
		Page string `json:"Page"`
	}

	c := colly.NewCollector()
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text // Get the text content of the <script> tag
		re := regexp.MustCompile(`vm.CurChapter = \{(.*?)\};`)
		baseImages = re.FindString(scriptContent)
	})

	if err := c.Visit(fmt.Sprintf("https://www.mangasee123.com/read-online/%v-chapter-%v.html", mangaId, chapterId)); err != nil {
		log.Info().Err(err)
		return returnData, err
	}

	baseImages = strings.Split(baseImages, "vm.CurChapter = ")[1]
	baseImages = strings.Split(baseImages, ";")[0]

	var mangaData MangaseeDataImage
	err := json.Unmarshal([]byte(baseImages), &mangaData)
	if err != nil {
		log.Info().Err(err)
		return returnData, err
	}

	count, err := strconv.Atoi(mangaData.Page)
	if err != nil {
		log.Info().Err(err)
		return returnData, err
	}

	for i := 1; i <= count; i++ {
		dataImages = append(dataImages, entity.Image{
			Image: fmt.Sprintf("https://official.lowee.us/manga/The-Final-Raid-Boss/%v-%v.png", chapterId, generateNumber(i)),
		})
	}

	returnData.ChapterName = chapterId
	returnData.Images = dataImages
	returnData.OriginalServer = fmt.Sprintf("https://www.mangasee123.com/read-online/%v-chapter-%v.html", mangaId, chapterId)
	return returnData, nil

}
