// main.go
package main

import (
	"context"
	"log"
	"msmc/src/browser"
	"msmc/src/config"
	"msmc/src/database"
	"msmc/src/jobs"
	"msmc/src/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
	browserDone := make(chan struct{})
	go func() {
		browser.Start(ctx, db)
		close(browserDone)
	}()

	// fiber
	app := server.NewServer(db, cfg)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		jobManager.Stop() // 크론 종료

		cancel()      // 브라우저 종료
		<-browserDone // 브라우저 종료 대기

		// 서버 종료
		if err := app.Shutdown(); err != nil {
			log.Printf("서버 종료 중 오류 발생: %v", err)
		}
	}()

	if err := app.Listen(cfg.Port); err != nil {
		log.Fatal(err)
	}
}
