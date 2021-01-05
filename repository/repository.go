package repository

import (
	"github.com/phpunch/route-roam-api/infrastructure/db"
	"github.com/phpunch/route-roam-api/model"
)

type Repository interface {
	fileRepository
	AddUser(user *model.User) error
	GetUser(email string) (*model.User, error)
	CreatePost(post *model.Post) error
	LikePost(like *model.Like) error
	UnlikePost(like *model.Like) error
	GetPosts() ([]model.Post, error)
	// CommentPost(userId string, postId string, text string) error
}

type repository struct {
	Ds *db.DB
}

func NewRepository(ds *db.DB) Repository {
	return &repository{
		Ds: ds,
	}
}

func (r *repository) AddUser(user *model.User) error {
	return r.Ds.PostgresqlDB.Insert(user)
}
func (r *repository) GetUser(email string) (*model.User, error) {
	return r.Ds.PostgresqlDB.QueryUser(email)
}
func (r *repository) CreatePost(post *model.Post) error {
	return r.Ds.PostgresqlDB.Insert(post)
}
func (r *repository) LikePost(like *model.Like) error {
	return r.Ds.PostgresqlDB.Insert(like)
}
func (r *repository) UnlikePost(like *model.Like) error {
	return r.Ds.PostgresqlDB.DeleteUserLike(like)
}
func (r *repository) GetPosts() ([]model.Post, error) {
	return r.Ds.PostgresqlDB.GetPosts()
}
