// Package database database/models.go
package database

import (
	"time"
)

type Post struct {
	ID        string    `gorm:"type:varchar;primaryKey"`
	PostID    int       `gorm:"column:postId;uniqueIndex:Post_postId_key"`
	Category  *string   `gorm:"type:varchar"`
	Title     string    `gorm:"type:varchar"`
	Author    string    `gorm:"type:varchar;index:Post_author_createdAt_idx"`
	CreatedAt time.Time `gorm:"column:createdAt;index:Post_author_createdAt_idx"`
}

func (Post) TableName() string {
	return "Post"
}
