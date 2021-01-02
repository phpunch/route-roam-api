package model

import (
	"github.com/lib/pq"
)

type Post struct {
	ID       int `gorm:"primaryKey"`
	UserID   string
	Text     *string
	ImageURL pq.StringArray `gorm:"type:text[]"`
}
