package postgresqlDB

import (
	"context"
	"fmt"
	"github.com/phpunch/route-roam-api/log"
	"github.com/phpunch/route-roam-api/model"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB interface {
	UserDBInterface
	// Upsert(data interface{}, clause clause.OnConflict) error
	// Insert(data interface{}) error
	// DeleteUserLike(like *model.Like) error
	GetPosts() ([]model.Post, error)

	Close() error
}

type PostgresqlDB struct {
	logger log.Logger
	DB     *pgxpool.Pool
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

	//db, err = pgx.Connect(context.Background(), connStr)
	var connectConf, _ = pgxpool.ParseConfig(connStr)
	connectConf.MaxConns = config.MaxOpenConns
	//connectConf.MaxConnLifetime = 300 * time.Second // use defaults until we have benchmarked this further
	//connectConf.HealthCheckPeriod = 300 * time.Second
	//connectConf.ConnConfig.PreferSimpleProtocol = true // don't wrap queries into transactions
	connectConf.ConnConfig.Logger = NewDatabaseLogger(&pgdb.logger)
	connectConf.ConnConfig.LogLevel = pgx.LogLevelWarn
	pgdb.DB, err = pgxpool.ConnectConfig(context.Background(), connectConf)
	if err != nil {
		pgdb.logger.Errorf("Error connecting to postgres: %+v")
		return nil, err
	}

	return pgdb, nil
}

func (pgdb *PostgresqlDB) Close() error {
	pgdb.DB.Close()
	return nil
}

// func (pgdb *PostgresqlDB) Upsert(data interface{}, clause clause.OnConflict) error {
// 	result := pgdb.DB.Clauses(clause).Create(data)
// 	err := result.Error
// 	if err != nil {
// 		return fmt.Errorf("upsert error: %v", err)
// 	}
// 	return nil
// }

// func (pgdb *PostgresqlDB) Insert(data interface{}) error {
// 	result := pgdb.DB.Create(data)
// 	err := result.Error
// 	if err != nil {
// 		return fmt.Errorf("insert error: %v", err)
// 	}
// 	return nil
// }

// func (pgdb *PostgresqlDB) DeleteUserLike(like *model.Like) error {
// 	result := pgdb.DB.Where("user_id = ? AND post_id = ?",
// 		like.UserID,
// 		like.PostID,
// 	).Delete(like)
// 	err := result.Error
// 	if err != nil {
// 		return fmt.Errorf("delete error: %v", err)
// 	}
// 	return nil
// }
