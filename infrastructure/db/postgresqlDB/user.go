package postgresqlDB

import (
	"context"
	"fmt"
	"github.com/phpunch/route-roam-api/model"
)

type UserDBInterface interface {
	CreateUser(user *model.User) (int64, error)
	QueryUser(username string) (*model.User, error)
}

func (pgdb *PostgresqlDB) CreateUser(user *model.User) (int64, error) {
	var userID int64
	err := pgdb.DB.QueryRow(context.Background(), `
		INSERT INTO users (
			"username",
			"password"
		)
		VALUES ($1, $2)
		RETURNING id
	`,
		user.Username,
		user.Password,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %v", err)
	}

	return userID, nil

}

func (pgdb *PostgresqlDB) QueryUser(username string) (*model.User, error) {
	var result model.User
	err := pgdb.DB.QueryRow(context.Background(), `
		SELECT id, username, password FROM users 
		WHERE users.username=$1
	`,
		username,
	).Scan(&result.ID, &result.Username, &result.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %v", err)
	}
	return &result, nil
}
