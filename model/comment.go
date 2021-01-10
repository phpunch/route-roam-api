package model

type Comment struct {
	ID     int64  `json:"id"`
	PostID int64  `json:"post_id"`
	UserID int64  `json:"user_id"`
	Text   string `json:"text"`
}
