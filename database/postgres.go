package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDb struct {
	PDb *pgxpool.Pool
}

const (
	host     = "localhost"
	port     = "5436"
	user     = "postgres"
	password = "12345"
	dbname   = "task_db"
)

func NewPostgresDB() (*PostgresDb, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		return nil, err
	}
	pgxPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return &PostgresDb{PDb: pgxPool}, nil
}
