package db

import (
	"context"
	"fmt"
	"github.com/phpunch/route-roam-api/infrastructure/db/minioDB"
	"github.com/phpunch/route-roam-api/infrastructure/db/postgresqlDB"
	"github.com/phpunch/route-roam-api/log"
)

type DB struct {
	PostgresqlDB postgresqlDB.DB
	MinioDB      miniodb.DB
}

func NewDB(logger log.Logger) (*DB, error) {

	postgresqlConfig, err := postgresqlDB.InitConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres config: %v", err)
	}
	postgresqlDB, err := postgresqlDB.New(postgresqlConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres database: %v", err)
	}

	minioConfig, err := miniodb.InitConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to init minio config: %v", err)
	}
	minioDB := miniodb.New(minioConfig)
	ctx := context.Background()
	if err := minioDB.CreateBucket(ctx, minioConfig.BucketName); err != nil {
		return nil, err
	}

	return &DB{
		PostgresqlDB: postgresqlDB,
		MinioDB:      minioDB,
	}, nil
}
