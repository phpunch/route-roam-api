package model

type Comment struct {
	ID       int64  `json:"id"`
	PostID   int64  `json:"post_id"`
	UserID   int64  `json:"user_id"`
	UserName string `json:"username"`
	Text     string `json:"text"`
}
