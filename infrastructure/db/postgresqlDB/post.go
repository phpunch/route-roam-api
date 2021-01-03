package postgresqlDB

// import (
// 	"github.com/phpunch/route-roam-api/model"
// )

// type PostDBInterface interface {
// 	LikePost(like model.Like) (*model.Like, error)
// }

// func (pgdb *PostgresqlDB) LikePost(like model.Like) (*model.Post, error) {
// 	var result model
// 	tx := pgdb.DB.Table("likes").Where("posts.likes = ?", email).First(&result)
// 	if tx.Error != nil {
// 		return nil, tx.Error
// 	}
// 	return &result, nil
// }
