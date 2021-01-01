package repository

import (
	"github.com/phpunch/route-roam-api/infrastructure/db"
	"github.com/phpunch/route-roam-api/model"
)

type Repository interface {
	AddUser(user *model.User) error
	GetUser(email string) (*model.User, error)
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
