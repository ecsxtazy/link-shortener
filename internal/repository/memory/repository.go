package memory

import (
	"context"
	"link-shortener/internal/model"
	"link-shortener/internal/repository"
	"sync"
)

type MemoryRepository struct {
	mu sync.RWMutex

	shortToLink map[string]model.Link
	urlToShort  map[string]string
}

func New() *MemoryRepository {
	return &MemoryRepository{
		shortToLink: make(map[string]model.Link),
		urlToShort:  make(map[string]string),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, link model.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.shortToLink[link.ShortCode] = link
	r.urlToShort[link.OriginalURL] = link.ShortCode

	return nil
}

func (r *MemoryRepository) GetbyShort(ctx context.Context, short string) (model.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	link, ok := r.shortToLink[short]
	if !ok {
		return model.Link{}, repository.ErrNotFound
	}
	return link, nil
}

func (r *MemoryRepository) GetByOriginal(ctx context.Context, url string) (model.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	short, ok := r.urlToShort[url]
	if !ok {
		return model.Link{}, repository.ErrNotFound
	}
	return r.shortToLink[short], nil
}
