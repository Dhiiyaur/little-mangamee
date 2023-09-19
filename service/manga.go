package service

import (
	"context"
	"little_mangamee/entity"
)

type MangaService interface {
	MangareadIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	MangareadSearch(ctx context.Context, search string) ([]entity.SearchData, error)
	MangareadDetail(ctx context.Context, mangaId string) (entity.DetailData, error)
	MangareadChapter(ctx context.Context, mangaId string) (entity.ChapterData, error)
	MangareadImage(ctx context.Context, mangaId string, chapterId string) (entity.ImageData, error)
	//
	MangabatIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	MangabatSearch(ctx context.Context, search string) ([]entity.SearchData, error)
	MangabatDetail(ctx context.Context, mangaId string) (entity.DetailData, error)
	MangabatChapter(ctx context.Context, mangaId string) (entity.ChapterData, error)
	MangabatImage(ctx context.Context, chapterId string) (entity.ImageData, error)
	//
	MangatownIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	MangatownSearch(ctx context.Context, search string) ([]entity.SearchData, error)
	MangatownDetail(ctx context.Context, mangaId string) (entity.DetailData, error)
	MangatownChapter(ctx context.Context, mangaId string) (entity.ChapterData, error)
	MangatownImage(ctx context.Context, mangaId string, chapterId string) (entity.ImageData, error)
	//
	MaidmyIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	MaidmySearch(ctx context.Context, search string) ([]entity.SearchData, error)
	MaidmyDetail(ctx context.Context, mangaId string) (entity.DetailData, error)
	MaidmyChapter(ctx context.Context, mangaId string) (entity.ChapterData, error)
	MaidmyImage(ctx context.Context, chapterId string) (entity.ImageData, error)
	//
	AsuraComicIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	AsuraComicSearch(ctx context.Context, search string) ([]entity.SearchData, error)
	AsuraComicDetail(ctx context.Context, mangaId string) (entity.DetailData, error)
	AsuraComicChapter(ctx context.Context, mangaId string) (entity.ChapterData, error)
	AsuraComicImage(ctx context.Context, chapterId string) (entity.ImageData, error)
	//
	ManganatoIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	ManganatoSearch(ctx context.Context, search string) ([]entity.SearchData, error)
	ManganatoDetail(ctx context.Context, mangaId string) (entity.DetailData, error)
	ManganatoChapter(ctx context.Context, mangaId string) (entity.ChapterData, error)
	ManganatoImage(ctx context.Context, mangaId string, chapterId string) (entity.ImageData, error)
	//
	ManganeloIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	ManganeloSearch(ctx context.Context, search string) ([]entity.SearchData, error)
	ManganeloDetail(ctx context.Context, mangaId string) (entity.DetailData, error)
	ManganeloChapter(ctx context.Context, mangaId string) (entity.ChapterData, error)
	ManganeloImage(ctx context.Context, mangaId string, chapterId string) (entity.ImageData, error)
}
