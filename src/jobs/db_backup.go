// jobs/backup.go
package jobs

import (
	"fmt"
	"log"
	"time"

	"msmc/src/util"
)

type BackupJob struct {
	dbPath     string
	bucketName string
}

func NewBackupJob(dbPath string) *BackupJob {
	return &BackupJob{
		dbPath:     dbPath,
		bucketName: "msmc",
	}
}

const (
	MAX_BACKUPS   = 5
	BACKUP_PREFIX = "db_backup/"
)

func (j *BackupJob) Execute() {
	client, err := util.CreateClient()
	if err != nil {
		log.Printf("S3 클라이언트 생성 실패: %v", err)
		return
	}

	now := time.Now()
	fileName := fmt.Sprintf("%smsmc_backup_%s.db",
		BACKUP_PREFIX,
		now.Format("2006-01-02_150405"),
	)

	err = util.UploadFile(client, j.bucketName, fileName, j.dbPath)
	if err != nil {
		log.Printf("백업 파일 업로드 실패: %v", err)
		return
	}

	files, err := util.ListFiles(client, j.bucketName, BACKUP_PREFIX)
	if err != nil {
		log.Printf("백업 파일 목록 조회 실패: %v", err)
		return
	}

	if len(files) > MAX_BACKUPS {
		filesToDelete := files[MAX_BACKUPS:]
		err = util.DeleteFiles(client, j.bucketName, filesToDelete)
		if err != nil {
			log.Printf("오래된 백업 파일 삭제 실패: %v", err)
			return
		}
	}

	log.Println("DB 백업 완료")
}
