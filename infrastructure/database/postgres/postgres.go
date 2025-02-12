package postgres

import (
	"context"
	"github.com/dhiemaz/fin-go/infrastructure/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDBConnection() (*pgxpool.Pool, error) {
	// Connection pool configuration
	dbURL := "postgres://username:password@localhost:5432/database_name"
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.WithFields(logger.Fields{"component": "infrastructure", "action": "init database connection"}).
			Infof("create database connection")
	}
	
	return pool, err
}
