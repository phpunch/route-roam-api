package postgresqlDB

import (
	"context"
	"fmt"
	"github.com/phpunch/route-roam-api/log"
	"github.com/phpunch/route-roam-api/model"
)

type PostDBInterface interface {
	CreatePost(post *model.Post) (int64, error)
	LikePost(like *model.Like) error
	UnlikePost(like *model.Like) error
	GetPosts() ([]model.Post, error)
	CreateComment(comment *model.Comment) error
	DeletePost(postID int64) error
	GetCommentsByPostID(postID int64) ([]model.Comment, error)
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
		log.Log.Debugf("%+v", err)

		return 0, fmt.Errorf("failed to create post")
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
		log.Log.Debugf("%+v", err)
		return fmt.Errorf("failed to like post")
	}

	return nil
}
func (pgdb *PostgresqlDB) UnlikePost(like *model.Like) error {
	commandTag, err := pgdb.DB.Exec(context.Background(), `
		delete from likes 
		where user_id=$1 and post_id=$2
`,
		like.UserID,
		like.PostID,
	)
	if err != nil {
		log.Log.Debugf("%+v", err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("failed to unlike post")
	}

	return nil
}

func (pgdb *PostgresqlDB) GetPosts() ([]model.Post, error) {
	var result []model.Post
	rows, err := pgdb.DB.Query(context.Background(), `
		select p.id, p.user_id, u.username, p.text, p.image_url, array_remove(array_agg(l.user_id), NULL) as liked_by, (
			select count(*) from comments com
			where com.post_id = p.id
		) as num_comments
		from posts p
		left join likes l on p.id = l.post_id
		left join users u on p.user_id = u.id
		group by p.id, u.username
	`)
	if err != nil {
		log.Log.Debugf("%+v", err)
		return nil, fmt.Errorf("failed to get posts")
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.UserID, &post.UserName, &post.Text, &post.ImageURL, &post.LikedBy, &post.NumComments)
		if err != nil {
			log.Log.Debugf("%+v", err)
			return nil, fmt.Errorf("failed to get posts")
		}
		result = append(result, post)
	}

	return result, nil
}

func (pgdb *PostgresqlDB) CreateComment(comment *model.Comment) error {
	err := pgdb.DB.QueryRow(context.Background(), `
		with rows as (INSERT INTO comments (
				"post_id",
				"user_id",
				"text"
			)
			VALUES ($1, $2, $3)
			RETURNING *
		)

		select r.id, u.username
		from rows r
		left join users u on r.user_id = u.id
`,
		comment.PostID,
		comment.UserID,
		comment.Text,
	).Scan(&comment.ID, &comment.UserName)
	if err != nil {
		log.Log.Debugf("%+v", err)
		return fmt.Errorf("failed to create comment")
	}

	return nil
}

func (pgdb *PostgresqlDB) DeletePost(postID int64) error {
	commandTag, err := pgdb.DB.Exec(context.Background(), `
		delete from posts where id=$1
	`, postID)
	if err != nil {
		log.Log.Debugf("%+v", err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("failed to delete post")
	}

	return nil
}

func (pgdb *PostgresqlDB) GetCommentsByPostID(postID int64) ([]model.Comment, error) {
	var result []model.Comment
	rows, err := pgdb.DB.Query(context.Background(), `
		select c.id, c.post_id, c.user_id, u.username, c.text 
		from comments c
		left join users u on c.user_id = u.id
		where c.post_id=$1
	`,
		postID,
	)
	if err != nil {
		log.Log.Debugf("%+v", err)
		return nil, fmt.Errorf("failed to query comment")
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Comment
		err = rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.UserName, &c.Text)
		if err != nil {
			log.Log.Debugf("%+v", err)
			return nil, fmt.Errorf("failed to query comment")
		}
		result = append(result, c)
	}
	return result, nil
}
