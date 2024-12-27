// Package jobs jobs/jobs.go
package jobs

import (
	"log"
	"msmc/src/browser"
	"msmc/src/config"
	"os"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type JobManager struct {
	cron *cron.Cron
	cfg  *config.Config
	db   *gorm.DB
}

func NewJobManager(cfg *config.Config, db *gorm.DB) *JobManager {
	return &JobManager{
		cron: cron.New(),
		cfg:  cfg,
		db:   db,
	}
}

func (jm *JobManager) StartJobs() {
	cleanupJob := NewArcaliveCleanupJob(jm.db, "https://arca.live/b/genshinskinmode/audit")
	backupJob := NewBackupJob(os.Getenv("DB_PATH"))

	// 2분마다 실행
	_, err := jm.cron.AddFunc("*/5 * * * *", cleanupJob.Execute)
	if err != nil {
		log.Printf("Failed to add cleanup job: %v", err)
		return
	}

	// 1시간마다 실행
	_, err = jm.cron.AddFunc("0 */1 * * *", func() {
		if err := browser.TriggerPageReload(); err != nil {
			log.Printf("Page reload failed: %v", err)
		}
	})
	if err != nil {
		log.Printf("Failed to add page reload job: %v", err)
		return
	}

	// 30분마다 실행
	_, err = jm.cron.AddFunc("*/30 * * * *", backupJob.Execute)
	if err != nil {
		log.Printf("Failed to add backup job: %v", err)
		return
	}

	jm.cron.Start()
	log.Println("All cron jobs started")
}

func (jm *JobManager) Stop() {
	<-jm.cron.Stop().Done()
	log.Println("All cron jobs stopped")
}
