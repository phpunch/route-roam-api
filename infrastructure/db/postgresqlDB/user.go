package postgresqlDB

import (
	"github.com/phpunch/route-roam-api/model"
)

type UserDBInterface interface {
	QueryUser(email string) (*model.User, error)
}

func (pgdb *PostgresqlDB) QueryUser(email string) (*model.User, error) {
	var result model.User
	tx := pgdb.DB.Table("users").Where("users.email = ?", email).First(&result)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &result, nil
}
