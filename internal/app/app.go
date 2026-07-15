package app

import (
	"context"
	"link-shortener/internal/config"
	"link-shortener/internal/middleware"
	"link-shortener/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg    config.Config
	server *http.Server
}

func New() *App {
	cfg := config.Load()
	mux := router.New()
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middleware.Logging(mux),
	}
	return &App{cfg: cfg, server: server}
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
	return app.server.Shutdown(ctx)
}
