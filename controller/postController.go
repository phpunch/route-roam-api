package controller

type postController interface {
	CreatePost(userId string, text string, image string) (bool, error)
	LikePost(userId string, postId string) (bool, error)
	CommentPost(userId string, postId string, text string) (bool, error)
}
