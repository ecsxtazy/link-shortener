package app

import (
	"context"
	"link-shortener/internal/config"
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
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
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
	log.Println("Stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return app.server.Shutdown(ctx)
}
