package model

import (
	"github.com/lib/pq"
)

type Post struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	UserID   string
	Text     *string
	ImageURL pq.StringArray `gorm:"type:text[]"`
	Likes    []User         `gorm:"many2many:likes"`
}
