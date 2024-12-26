package arcalive

import (
	"errors"
	"fmt"
	"log"
	"msmc/src/database"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Crawler struct {
	db *gorm.DB
}

func NewCrawler(db *gorm.DB) *Crawler {
	return &Crawler{
		db: db,
	}
}

// 첫 페이지
func (c *Crawler) GetPost(url string) (int, int, error) {
	fmt.Println("현재 처리중인 URL:", url)

	result, err := GetPosts(url, 1, "")
	if err != nil {
		return 0, 0, fmt.Errorf("데이터 가져오기 실패: %v", err)
	}

	return c.processPosts(result.Posts)
}

// 모든 페이지
func (c *Crawler) GetPosts(url string) (int, int, error) {
	currentURL := url
	var totalCreated, totalUpdated int

	for currentURL != "" {
		fmt.Println("현재 처리중인 URL:", currentURL)

		result, err := GetPosts(currentURL, 1, "")
		if err != nil {
			return totalCreated, totalUpdated, fmt.Errorf("데이터 가져오기 실패: %v", err)
		}

		created, updated, err := c.processPosts(result.Posts)
		if err != nil {
			return totalCreated, totalUpdated, err
		}

		totalCreated += created
		totalUpdated += updated
		fmt.Printf("현재 페이지: %d개 중 %d개 생성, %d개 업데이트됨\n",
			len(result.Posts), created, updated)

		// 다음 페이지 URL 설정
		currentURL = result.NextPageURL

		// 페이지 간 딜레이
		if currentURL != "" {
			time.Sleep(3 * time.Second)
		}
	}

	return totalCreated, totalUpdated, nil
}

func (c *Crawler) processPosts(posts []Post) (int, int, error) {
	createdCount, updatedCount := 0, 0

	for _, post := range posts {
		var existing database.Post
		result := c.db.First(&existing, "postId = ?", post.PostID)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 새로운 포스트 생성
			id := uuid.New()
			newPost := database.Post{
				ID:        id.String(),
				PostID:    post.PostID,
				Category:  &post.Category,
				Title:     post.Title,
				Author:    post.Author,
				CreatedAt: post.CreatedAt,
			}

			if err := c.db.Create(&newPost).Error; err != nil {
				log.Printf("포스트 생성 실패 (postId: %d): %v", post.PostID, err)
				continue
			}
			createdCount++
		} else if result.Error == nil {
			// 데이터가 다른 경우에만 업데이트
			if existing.Category != nil && *existing.Category != post.Category ||
				existing.Title != post.Title ||
				existing.Author != post.Author ||
				existing.CreatedAt.Unix() != post.CreatedAt.Unix() {

				existing.Category = &post.Category
				existing.Title = post.Title
				existing.Author = post.Author
				existing.CreatedAt = post.CreatedAt

				if err := c.db.Save(&existing).Error; err != nil {
					log.Printf("포스트 업데이트 실패 (postId: %d): %v", post.PostID, err)
					continue
				}
				updatedCount++
			}
		} else {
			log.Printf("포스트 조회 실패 (postId: %d): %v", post.PostID, result.Error)
			continue
		}
	}

	return createdCount, updatedCount, nil
}
