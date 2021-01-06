package model

type Post struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	UserID   int
	Text     *string
	ImageURL []string `gorm:"type:text[]"`
	Likes    int
}
