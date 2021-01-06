package postgresqlDB

import (
	"context"
	"fmt"
	"github.com/phpunch/route-roam-api/model"
)

type PostDBInterface interface {
	CreatePost(post *model.Post) error
	LikePost(like *model.Like) error
	UnlikePost(like *model.Like) error
	GetPosts() ([]model.Post, error)
}

func (pgdb *PostgresqlDB) CreatePost(post *model.Post) error {
	_, err := pgdb.DB.Exec(context.Background(), `
	INSERT INTO posts (
		"user_id",
		"text",
		"image_url"
	)
	VALUES ($1, $2, $3)
`,
		post.UserID,
		post.Text,
		post.ImageURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create post: %v", err)
	}

	return nil

}
func (pgdb *PostgresqlDB) LikePost(like *model.Like) error {
	return nil
}
func (pgdb *PostgresqlDB) UnlikePost(like *model.Like) error {
	return nil
}

func (pgdb *PostgresqlDB) GetPosts() ([]model.Post, error) {
	var result []model.Post
	rows, err := pgdb.DB.Query(context.Background(), `
		SELECT id, user_id, text, image_url FROM posts 
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.UserID, &post.Text, &post.ImageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		result = append(result, post)
	}
	return result, nil
}
