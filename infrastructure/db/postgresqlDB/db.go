package postgresqlDB

import (
	"fmt"
	"github.com/phpunch/route-roam-api/log"
	"github.com/phpunch/route-roam-api/model"

	"gorm.io/gorm/clause"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB interface {
	UserDBInterface
	Upsert(data interface{}, clause clause.OnConflict) error
	Insert(data interface{}) error
	DeleteUserLike(like *model.Like) error
	GetPosts() ([]model.Post, error)

	Close() error
}

type PostgresqlDB struct {
	logger log.Logger
	DB     *gorm.DB
}

func New(config *Config, logger log.Logger) (pgdb *PostgresqlDB, err error) {
	pgdb = &PostgresqlDB{
		logger: logger.WithFields(log.Fields{
			"module": "db",
		}),
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUsername,
		config.DBPassword,
		config.DBName,
	)
	dialector := postgres.Open(connStr)

	pgdb.DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: NewGormLogger(&pgdb.logger),
	})
	if err != nil {
		pgdb.logger.Errorf("Error connecting to postgres: %+v", err)
		return nil, err
	}

	if err := pgdb.DB.SetupJoinTable(&model.Post{}, "Likes", &model.Like{}); err != nil {
		return nil, err
	}

	if err := pgdb.DB.AutoMigrate(&model.Post{}, &model.User{}); err != nil {
		return nil, err
	}

	return pgdb, nil
}

func (pgdb *PostgresqlDB) Close() error {
	connDB, err := pgdb.DB.DB()
	if err != nil {
		pgdb.logger.Errorf("Errorf close db: %+v", err)
		return err
	}
	err = connDB.Close()
	return err
}

func (pgdb *PostgresqlDB) Upsert(data interface{}, clause clause.OnConflict) error {
	result := pgdb.DB.Clauses(clause).Create(data)
	err := result.Error
	if err != nil {
		return fmt.Errorf("upsert error: %v", err)
	}
	return nil
}

func (pgdb *PostgresqlDB) Insert(data interface{}) error {
	result := pgdb.DB.Create(data)
	err := result.Error
	if err != nil {
		return fmt.Errorf("insert error: %v", err)
	}
	return nil
}

func (pgdb *PostgresqlDB) DeleteUserLike(like *model.Like) error {
	result := pgdb.DB.Where("user_id = ? AND post_id = ?",
		like.UserID,
		like.PostID,
	).Delete(like)
	err := result.Error
	if err != nil {
		return fmt.Errorf("delete error: %v", err)
	}
	return nil
}
