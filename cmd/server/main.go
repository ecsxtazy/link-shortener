package main

import (
	"link-shortener/internal/app"
	"log"
)

func main() {
	app := app.New()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
