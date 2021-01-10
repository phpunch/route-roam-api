package postgresqlDB

import (
	"context"
	"fmt"
	"github.com/phpunch/route-roam-api/model"
)

type PostDBInterface interface {
	CreatePost(post *model.Post) (int64, error)
	LikePost(like *model.Like) error
	UnlikePost(like *model.Like) error
	GetPosts() ([]model.Post, error)
}

func (pgdb *PostgresqlDB) CreatePost(post *model.Post) (int64, error) {
	var postID int64
	err := pgdb.DB.QueryRow(context.Background(), `
	INSERT INTO posts (
		"user_id",
		"text",
		"image_url"
	)
	VALUES ($1, $2, $3)
	RETURNING id
`,
		post.UserID,
		post.Text,
		post.ImageURL,
	).Scan(&postID)
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %v", err)
	}

	return postID, nil
}
func (pgdb *PostgresqlDB) LikePost(like *model.Like) error {
	_, err := pgdb.DB.Exec(context.Background(), `
	INSERT INTO likes (
		"user_id",
		"post_id"
	)
	VALUES ($1, $2)
`,
		like.UserID,
		like.PostID,
	)
	if err != nil {
		return fmt.Errorf("failed to create post: %v", err)
	}

	return nil
}
func (pgdb *PostgresqlDB) UnlikePost(like *model.Like) error {
	commandTag, _ := pgdb.DB.Exec(context.Background(), `
		delete from likes 
		where user_id=$1 and post_id=$2
`,
		like.UserID,
		like.PostID,
	)
	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("failed to delete row")
	}

	return nil
}

func (pgdb *PostgresqlDB) GetPosts() ([]model.Post, error) {
	var result []model.Post
	rows, err := pgdb.DB.Query(context.Background(), `
		select p.id, p.user_id, p.text, p.image_url, array_agg(l.user_id) as liked_by from posts p
		join likes l on p.id = l.post_id
		group by p.id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.UserID, &post.Text, &post.ImageURL, &post.LikedBy)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		result = append(result, post)
	}
	return result, nil
}
