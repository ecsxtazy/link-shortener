package memory

import (
	"context"
	"fmt"
	"link-shortener/internal/model"
	"link-shortener/internal/repository"
	"link-shortener/internal/repository/memory"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	repo := memory.New()
	link := model.Link{
		ShortCode:   "abcdefghij",
		OriginalURL: "https://google.com",
	}
	err := repo.Create(context.Background(), link)
	require.NoError(t, err)
	got, err := repo.GetByShort(context.Background(), "abcdefghij")
	require.NoError(t, err)
	require.Equal(t, link, got)
}

func TestGetByShort_NotFound(t *testing.T) {
	repo := memory.New()
	_, err := repo.GetByShort(context.Background(), "")
	require.ErrorIs(t, err, repository.ErrNotFound)
}

func TestGetByURL_NotFound(t *testing.T) {
	repo := memory.New()
	_, err := repo.GetByOriginal(context.Background(), "")
	require.ErrorIs(t, err, repository.ErrNotFound)
}

func TestRepository_ConcurrentCreate(t *testing.T) {
	repo := memory.New()
	const workers = 100
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			link := model.Link{
				ShortCode:   fmt.Sprintf("code%06d", i),
				OriginalURL: fmt.Sprintf("https://site%d.com", i),
			}
			err := repo.Create(context.Background(), link)
			require.NoError(t, err)

		}(i)
	}
	wg.Wait()
	for i := 0; i < workers; i++ {
		link, err := repo.GetByShort(
			context.Background(),
			fmt.Sprintf("code%06d", i),
		)
		require.NoError(t, err)
		require.Equal(
			t,
			fmt.Sprintf("https://site%d.com", i),
			link.OriginalURL,
		)
	}
}
