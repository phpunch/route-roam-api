package db

import (
	"fmt"
	"github.com/phpunch/route-roam-api/infrastructure/db/postgresqlDB"
	"github.com/phpunch/route-roam-api/log"
)

type DB struct {
	PostgresqlDB postgresqlDB.DB
}

func NewDB(logger log.Logger) (*DB, error) {

	dbConfig, err := postgresqlDB.InitConfig()
	if err != nil {
		// TODO: error handling
		return nil, fmt.Errorf("failed to init postgres config: %v", err)
	}

	postgresqlDB, err := postgresqlDB.New(dbConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres database: %v", err)
	}
	return &DB{
		PostgresqlDB: postgresqlDB,
	}, nil
}
