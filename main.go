// main.go
package main

import (
	"log"
	"msmc/src/browser"
	"msmc/src/config"
	"msmc/src/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.LoadConfig()

	go browser.Start()

	app := server.NewServer(cfg)
	log.Fatal(app.Listen(cfg.Port))
}
