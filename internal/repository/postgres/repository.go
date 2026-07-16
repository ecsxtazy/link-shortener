package postgres

import (
	"context"
	"errors"
	"link-shortener/internal/model"
	"link-shortener/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) (*PostgresRepository, error) {
	r := &PostgresRepository{
		pool: pool,
	}
	if err := r.init(context.Background()); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *PostgresRepository) init(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS links (short_code CHAR(10) PRIMARY KEY,original_url TEXT UNIQUE NOT NULL);`
	_, err := r.pool.Exec(ctx, query)
	return err
}

func (r *PostgresRepository) Create(ctx context.Context, link model.Link) error {
	query := `INSERT INTO links(short_code, original_url) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, link.ShortCode, link.OriginalURL)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return repository.ErrAlreadyExists
		}
	}
	return err
}

func (r *PostgresRepository) GetByShort(ctx context.Context, short string) (model.Link, error) {
	query := `SELECT short_code, original_url FROM links WHERE short_code = $1`
	var link model.Link
	err := r.pool.QueryRow(ctx, query, short).Scan(&link.ShortCode, &link.OriginalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Link{}, repository.ErrNotFound
		}
		return model.Link{}, err
	}
	return link, nil
}

func (r *PostgresRepository) GetByOriginal(ctx context.Context, url string) (model.Link, error) {
	query := `SELECT short_code, original_url FROM links WHERE original_url = $1`
	var link model.Link
	err := r.pool.QueryRow(ctx, query, url).Scan(&link.ShortCode, &link.OriginalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Link{}, repository.ErrNotFound
		}
		return model.Link{}, err
	}
	return link, nil
}

func (r *PostgresRepository) Close() {
	r.pool.Close()
}
