package service

import (
	"context"
	"little_mangamee/entity"
)

type MangaService interface {
	MangareadSearch(ctx context.Context, search string) ([]entity.SearchData, error)
	MangatownSearch(ctx context.Context, search string) ([]entity.SearchData, error)
	MangabatSearch(ctx context.Context, search string) ([]entity.SearchData, error)
	MaidmySearch(ctx context.Context, search string) ([]entity.SearchData, error)
	//
	MangabatIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	MangareadIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	MangatownIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	MaidmyIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	//
	MangabatDetail(ctx context.Context, mangaId string) (entity.DetailData, error)
	MangareadDetail(ctx context.Context, mangaId string) (entity.DetailData, error)
	MangatownDetail(ctx context.Context, mangaId string) (entity.DetailData, error)
	MaidmyDetail(ctx context.Context, mangaId string) (entity.DetailData, error)
	//
	MangabatChapter(ctx context.Context, mangaId string) (entity.ChapterData, error)
	MangareadChapter(ctx context.Context, mangaId string) (entity.ChapterData, error)
	MangatownChapter(ctx context.Context, mangaId string) (entity.ChapterData, error)
	MaidmyChapter(ctx context.Context, mangaId string) (entity.ChapterData, error)
	//
	MaidmyImage(ctx context.Context, chapterId string) (entity.ImageData, error)
	MangatownImage(ctx context.Context, mangaId string, chapterId string) (entity.ImageData, error)
	MangareadImage(ctx context.Context, mangaId string, chapterId string) (entity.ImageData, error)
	MangabatImage(ctx context.Context, chapterId string) (entity.ImageData, error)
}
