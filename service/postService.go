package service

import (
	"github.com/phpunch/route-roam-api/model"
)

type postService interface {
	CreatePost(userId int64, text string, imageURLs []string) (*model.Post, error)
	LikePost(userId int64, postId int64) error
	UnlikePost(userId int64, postId int64) error
	CommentPost(userId int64, postId int64, text string) error
	GetPosts() ([]model.Post, error)
}

func (s *service) CreatePost(userId int64, text string, imageURLs []string) (*model.Post, error) {
	post := &model.Post{
		UserID:   userId,
		Text:     text,
		ImageURL: imageURLs,
	}
	postID, err := s.repository.CreatePost(post)
	if err != nil {
		return nil, err
	}
	post.ID = postID
	return post, nil
}
func (s *service) LikePost(userId int64, postId int64) error {
	like := &model.Like{
		UserID: userId,
		PostID: postId,
	}
	return s.repository.LikePost(like)
}
func (s *service) UnlikePost(userId int64, postId int64) error {
	like := &model.Like{
		UserID: userId,
		PostID: postId,
	}
	return s.repository.UnlikePost(like)
}
func (s *service) GetPosts() ([]model.Post, error) {
	return s.repository.GetPosts()
}

func (s *service) CommentPost(userId int64, postId int64, text string) error {
	return nil
}
