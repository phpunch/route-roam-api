package postgresqlDB

import (
	"github.com/phpunch/route-roam-api/model"
)

type PostDBInterface interface {
	GetPosts() ([]model.Post, error)
}

func (pgdb *PostgresqlDB) GetPosts() ([]model.Post, error) {
	var result []model.Post
	tx := pgdb.DB.Preload("likes").Find(&result)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return result, nil
}
