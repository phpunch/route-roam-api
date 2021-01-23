package model

type Post struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	UserName    string    `json:"username"`
	Text        string    `json:"text"`
	ImageURL    []string  `json:"image_url"`
	LikedBy     []*int    `json:"liked_by"`
	Comments    []Comment `json:"comments"`
	NumComments int64     `json:"num_comments"`
}
