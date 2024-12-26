// Package jobs jobs/arcalive_cleanup.go
package jobs

import (
	"gorm.io/gorm"
	"log"
	"msmc/src/arcalive"
	"msmc/src/database"
)

type ArcaliveCleanupJob struct {
	db      *gorm.DB
	baseURL string
}

func NewArcaliveCleanupJob(db *gorm.DB, baseURL string) *ArcaliveCleanupJob {
	return &ArcaliveCleanupJob{
		db:      db,
		baseURL: baseURL,
	}
}

func (j *ArcaliveCleanupJob) Execute() {
	audits, err := arcalive.GetAudits(j.baseURL)
	if err != nil {
		log.Printf("Failed to fetch audits: %v", err)
		return
	}

	var deletedPostIDs []int
	for _, audit := range audits {
		if audit.Action == arcalive.AuditActionDelete {
			deletedPostIDs = append(deletedPostIDs, audit.PostID)
		}
	}

	if len(deletedPostIDs) == 0 {
		return
	}

	result := j.db.Where("postId IN ?", deletedPostIDs).Delete(&database.Post{})
	if result.Error != nil {
		log.Printf("posts 삭제 실패: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("%d 개 글이 db에서 삭제됨", result.RowsAffected)
	}
}
