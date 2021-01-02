package service

type postService interface {
	CreatePost(userId string, text *string, image *string) error
	LikePost(userId string, postId string) error
	CommentPost(userId string, postId string, text string) error
}

func (s *service) CreatePost(userId string, text *string, image *string) error {
	return nil
}
func (s *service) LikePost(userId string, postId string) error {
	return nil
}
func (s *service) CommentPost(userId string, postId string, text string) error {
	return nil
}
