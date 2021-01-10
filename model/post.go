package model

type Post struct {
	ID       int64
	UserID   int64
	Text     string
	ImageURL []string
	LikedBy  []int
}
