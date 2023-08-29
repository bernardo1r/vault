package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func New(connString string) (DB, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	return DB{
		pool: pool,
	}, err
}
