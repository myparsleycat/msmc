// Package jobs jobs/jobs.go
package jobs

import (
	"log"
	"msmc/src/browser"
	"msmc/src/config"

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

	// 2분마다 실행
	_, err := jm.cron.AddFunc("*/2 * * * *", cleanupJob.Execute)
	if err != nil {
		log.Printf("Failed to add cleanup job: %v", err)
		return
	}

	// 1시간마다 실행
	_, err = jm.cron.AddFunc("* */1 * * *", browser.TriggerPageReload)
	if err != nil {
		log.Printf("Failed to add page reload job: %v", err)
		return
	}

	jm.cron.Start()
	log.Println("All cron jobs started")
}

func (jm *JobManager) Stop() {
	<-jm.cron.Stop().Done()
	log.Println("All cron jobs stopped")
}
