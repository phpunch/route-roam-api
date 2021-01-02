package model

type Post struct {
	ID       int     `gorm:"primaryKey"`
	UserID   string  `json:"userId" form:"userId"`
	Text     *string `json:"text" form:"text"`
	ImageURL *string `json:"imageUrl" form:"imageUrl"`
}
