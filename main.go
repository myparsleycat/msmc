// main.go
package main

import (
	"log"
	"msmc/src/browser"
	"msmc/src/config"
	"msmc/src/database"
	"msmc/src/jobs"
	"msmc/src/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.LoadConfig()
	db := database.GetDB()

	// cron
	jobManager := jobs.NewJobManager(cfg, db)
	jobManager.StartJobs()

	// chromedp
	go browser.Start(db)

	// fiber
	app := server.NewServer(db, cfg)
	log.Fatal(app.Listen(cfg.Port))
}
