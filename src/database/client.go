// Package database src/database/client.go
package database

import (
	"log"
	"os"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	instance *gorm.DB
	once     sync.Once
)

// GetDB 싱글톤
func GetDB() *gorm.DB {
	db_path := os.Getenv("DB_PATH") + "?cache=shared&_journal_mode=WAL"

	once.Do(func() {
		db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
		if err != nil {
			log.Fatal("데이터베이스 연결 실패:", err)
		}
		db.Config.Logger = logger.Default.LogMode(logger.Silent)
		instance = db
	})
	return instance
}

// InitModels 모델 초기화
func InitModels() error {
	db := GetDB()

	// 테이블이 이미 존재하는지 확인
	if db.Migrator().HasTable(&Post{}) {
		log.Println("Post 테이블이 이미 존재합니다. 스키마 검증만 수행합니다.")
		return nil
	}

	db = db.Set("gorm:auto_migrate_options", "SKIP_INDEX")

	// 테이블이 없는 경우에만 생성
	return db.AutoMigrate(&Post{})
}
