package model

type Post struct {
	ID       int64    `json:"id"`
	UserID   int64    `json:"user_id"`
	Text     string   `json:"text"`
	ImageURL []string `json:"image_url"`
	LikedBy  []int    `json:"liked_by"`
}
