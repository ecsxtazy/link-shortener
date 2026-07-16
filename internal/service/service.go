package service

import (
	"context"
	"errors"
	"link-shortener/internal/model"
	"link-shortener/internal/repository"
)

const maxGenerateAttempts = 5

type Shortener interface {
	Shorten(ctx context.Context, originalURL string) (string, error)
	Resolve(ctx context.Context, shortCode string) (string, error)
}

type Service struct {
	repo      repository.Repository
	generator Generator
}

func New(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Resolve(ctx context.Context, short string) (string, error) {
	link, err := s.repo.GetByShort(ctx, short)
	if err != nil {
		return "", err
	}
	return link.OriginalURL, nil
}

func (s *Service) Shorten(ctx context.Context, originalURL string) (string, error) {
	link, err := s.repo.GetByOriginal(ctx, originalURL)
	if err == nil {
		return link.ShortCode, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return "", err
	}
	for i := 0; i < maxGenerateAttempts; i++ {
		code, err := Generate()
		if err != nil {
			return "", err
		}
		link := model.Link{
			ShortCode:   code,
			OriginalURL: originalURL,
		}
		err = s.repo.Create(ctx, link)
		if err == nil {
			return code, nil
		}
		switch {
		case errors.Is(err, repository.ErrShortCodeExists):
			continue
		case errors.Is(err, repository.ErrOriginalURLExists):
			existing, err := s.repo.GetByOriginal(ctx, originalURL)
			if err != nil {
				return "", err
			}
			return existing.ShortCode, nil
		default:
			return "", err
		}
	}
	return "", ErrGenerationFailed
}
