package service

import (
	"context"
	"little_mangamee/entity"
)

type MangaService interface {
	MangabatIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	MangareadIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	MangatownIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
	MaidmyIndex(ctx context.Context, pageNumber string) ([]entity.IndexData, error)
}
