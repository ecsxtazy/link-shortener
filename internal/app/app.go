package app

import (
	"context"
	"fmt"
	"link-shortener/internal/config"
	"link-shortener/internal/middleware"
	"link-shortener/internal/repository"
	"link-shortener/internal/repository/memory"
	"link-shortener/internal/repository/postgres"
	"link-shortener/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg    config.Config
	server *http.Server
	repo   repository.Repository
}

func New() (*App, error) {
	cfg := config.Load()
	var repo repository.Repository
	switch cfg.Storage {
	case "memory":
		repo = memory.New()
	case "postgres":
		pool, err := pgxpool.New(context.Background(), cfg.DBUrl)
		if err != nil {
			return nil, err
		}
		repo, err = postgres.New(pool)
		if err != nil {
			pool.Close()
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown storage: %s", cfg.Storage)
	}
	mux := router.New()
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.Logging(mux),
	}
	return &App{cfg: cfg, server: server, repo: repo}, nil
}

func (app *App) Run() error {
	go func() {
		log.Printf("Server started on :%s", app.cfg.Port)
		if err := app.server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Stopping server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer app.repo.Close()
	return app.server.Shutdown(ctx)
}
