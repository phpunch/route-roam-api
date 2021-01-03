package service

import (
	"github.com/phpunch/route-roam-api/model"
)

type postService interface {
	CreatePost(userId string, text *string, imageURLs []string) error
	LikePost(userId string, postId string) error
	CommentPost(userId string, postId string, text string) error
}

func (s *service) CreatePost(userId string, text *string, imageURLs []string) error {
	post := &model.Post{
		UserID:   userId,
		Text:     text,
		ImageURL: imageURLs,
	}
	return s.repository.CreatePost(post)
}
func (s *service) LikePost(userId string, postId string) error {
	like := &model.Like{
		UserID: userId,
		PostID: postId,
	}
	return s.repository.LikePost(like)
}
func (s *service) CommentPost(userId string, postId string, text string) error {
	return nil
}
