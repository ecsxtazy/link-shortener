package repository

import (
	"context"
	"link-shortener/internal/model"
)

type Repository interface {
	Save(ctx context.Context, link model.Link) error
	GetByShort(ctx context.Context, short string) (model.Link, error)
	GetByOriginal(ctx context.Context, url string) (model.Link, error)
}
