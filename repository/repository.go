package repository

import (
	"github.com/phpunch/route-roam-api/infrastructure/db"
	"github.com/phpunch/route-roam-api/model"
	"time"
)

type Repository interface {
	fileRepository
	AddUser(user *model.User) (int64, error)
	SaveToken(uuid string, userID string, expDuration time.Duration) error
	FetchToken(uuid string) (int64, error)
	DeleteToken(uuid string) (int64, error)
	GetUser(email string) (*model.User, error)
	CreatePost(post *model.Post) (int64, error)
	LikePost(like *model.Like) error
	UnlikePost(like *model.Like) error
	GetPosts() ([]model.Post, error)
	CommentPost(comment *model.Comment) (int64, error)
	GetCommentsByPostID(postID int64) ([]model.Comment, error)
}

type repository struct {
	Ds *db.DB
}

func NewRepository(ds *db.DB) Repository {
	return &repository{
		Ds: ds,
	}
}

func (r *repository) AddUser(user *model.User) (int64, error) {
	return r.Ds.PostgresqlDB.CreateUser(user)
}
func (r *repository) SaveToken(uuid string, userID string, expDuration time.Duration) error {
	return r.Ds.RedisDB.Set(uuid, userID, expDuration)
}
func (r *repository) FetchToken(uuid string) (int64, error) {
	return r.Ds.RedisDB.Get(uuid)
}
func (r *repository) DeleteToken(uuid string) (int64, error) {
	return r.Ds.RedisDB.Del(uuid)
}
func (r *repository) GetUser(email string) (*model.User, error) {
	return r.Ds.PostgresqlDB.QueryUser(email)
}
func (r *repository) CreatePost(post *model.Post) (int64, error) {
	return r.Ds.PostgresqlDB.CreatePost(post)
}
func (r *repository) LikePost(like *model.Like) error {
	return r.Ds.PostgresqlDB.LikePost(like)
}
func (r *repository) UnlikePost(like *model.Like) error {
	return r.Ds.PostgresqlDB.UnlikePost(like)
}
func (r *repository) GetPosts() ([]model.Post, error) {
	return r.Ds.PostgresqlDB.GetPosts()
}
func (r *repository) CommentPost(comment *model.Comment) (int64, error) {
	return r.Ds.PostgresqlDB.CreateComment(comment)
}
func (r *repository) GetCommentsByPostID(postID int64) ([]model.Comment, error) {
	return r.Ds.PostgresqlDB.GetCommentsByPostID(postID)
}
