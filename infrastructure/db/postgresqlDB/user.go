package postgresqlDB

import (
	"context"
	"fmt"
	"github.com/phpunch/route-roam-api/model"
)

type UserDBInterface interface {
	CreateUser(user *model.User) error
	QueryUser(email string) (*model.User, error)
}

func (pgdb *PostgresqlDB) CreateUser(user *model.User) error {
	_, err := pgdb.DB.Exec(context.Background(), `
		INSERT INTO users (
			"email",
			"password"
		)
		VALUES ($1, $2)
	`,
		user.Email,
		user.Password,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil

}

func (pgdb *PostgresqlDB) QueryUser(email string) (*model.User, error) {
	var result model.User
	err := pgdb.DB.QueryRow(context.Background(), `
		SELECT id, email, password FROM users 
		WHERE users.email=$1
	`,
		email,
	).Scan(&result.ID, &result.Email, &result.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %v", err)
	}
	return &result, nil
}
