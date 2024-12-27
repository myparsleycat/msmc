package handlers

import (
	"fmt"
	"log"
	"msmc/src/browser"
	"msmc/src/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UtilsReloadHandler(c *fiber.Ctx, db *gorm.DB) error {
	err := browser.TriggerPageReload()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "페이지 새로고침 성공",
	})
}

const (
	MAX_BACKUPS   = 5
	BACKUP_PREFIX = "db_backup/"
)

func UtilsTestHandler(c *fiber.Ctx, db *gorm.DB) error {
	client, err := util.CreateClient()
	if err != nil {
		log.Printf("S3 클라이언트 생성 실패: %v", err)
	}

	now := time.Now()
	fileName := fmt.Sprintf("%smsmc_backup_%s.db",
		BACKUP_PREFIX,
		now.Format("2006-01-02_150405"),
	)

	err = util.UploadFile(client, "msmc", fileName, "msmc.db")
	if err != nil {
		log.Printf("백업 파일 업로드 실패: %v", err)
	}

	files, err := util.ListFiles(client, "msmc", BACKUP_PREFIX)
	if err != nil {
		log.Printf("백업 파일 목록 조회 실패: %v", err)
	}

	if len(files) > MAX_BACKUPS {
		filesToDelete := files[MAX_BACKUPS:]
		err = util.DeleteFiles(client, "msmc", filesToDelete)
		if err != nil {
			log.Printf("오래된 백업 파일 삭제 실패: %v", err)
		}
	}

	log.Println("DB 백업 완료")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"test": "test",
	})
}
