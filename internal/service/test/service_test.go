package test

import (
	"context"
	"errors"
	"link-shortener/internal/model"
	"link-shortener/internal/repository"
	"link-shortener/internal/repository/memory"
	"link-shortener/internal/service"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeGenerator struct {
	codes []string
	index int
	err   error
}

func (g *fakeGenerator) Generate() (string, error) {
	if g.err != nil {
		return "", g.err
	}
	code := g.codes[g.index]
	g.index++
	return code, nil
}

func TestShorten_NewURL(t *testing.T) {
	repo := memory.New()
	gen := &fakeGenerator{
		codes: []string{
			"abcdefghij",
		},
	}
	svc := service.New(repo, gen)
	code, err := svc.Shorten(
		context.Background(),
		"https://google.com",
	)
	require.NoError(t, err)
	require.Equal(t, "abcdefghij", code)
	link, err := repo.GetByOriginal(
		context.Background(),
		"https://google.com",
	)
	require.NoError(t, err)
	require.Equal(t, "abcdefghij", link.ShortCode)
}

func TestShorten_ExistingURL(t *testing.T) {
	repo := memory.New()
	err := repo.Create(
		context.Background(),
		model.Link{
			ShortCode:   "abcdefghij",
			OriginalURL: "https://google.com",
		},
	)
	require.NoError(t, err)
	gen := &fakeGenerator{}
	svc := service.New(repo, gen)
	code, err := svc.Shorten(
		context.Background(),
		"https://google.com",
	)
	require.NoError(t, err)
	require.Equal(t, "abcdefghij", code)
	require.Equal(t, 0, gen.index)
}

func TestShorten_GeneratorError(t *testing.T) {
	repo := memory.New()
	gen := &fakeGenerator{
		err: errors.New("boom"),
	}
	svc := service.New(repo, gen)
	_, err := svc.Shorten(
		context.Background(),
		"https://google.com",
	)
	require.Error(t, err)
	require.EqualError(t, err, "boom")
}

func TestResolve(t *testing.T) {
	repo := memory.New()
	err := repo.Create(
		context.Background(),
		model.Link{
			ShortCode:   "abcdefghij",
			OriginalURL: "https://google.com",
		},
	)
	require.NoError(t, err)
	svc := service.New(
		repo,
		&fakeGenerator{},
	)
	url, err := svc.Resolve(
		context.Background(),
		"abcdefghij",
	)
	require.NoError(t, err)
	require.Equal(
		t,
		"https://google.com",
		url,
	)
}

func TestResolve_NotFound(t *testing.T) {
	repo := memory.New()
	svc := service.New(
		repo,
		&fakeGenerator{},
	)
	_, err := svc.Resolve(
		context.Background(),
		"abcdefghij",
	)
	require.ErrorIs(
		t,
		err,
		repository.ErrNotFound,
	)
}
