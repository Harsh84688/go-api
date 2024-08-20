package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool
var Ctx = context.Background()

func InitDB() {
	// var err error
	// // Define connection string
	// connString := "postgres://user:password@db:5432/mydatabase"

	// // Connect to the database which will be created thru docker-compose
	// DB, err = pgx.Connect(Ctx, connString)
	// if err != nil {
	// 	panic("Unable to connect to database")
	// }

	config, err := pgxpool.ParseConfig("postgres://user:password@db:5432/mydatabase")
	if err != nil {
		panic("Unable to parse database")
	}

	// Configure connection pool settings
	config.MaxConns = 10                      // Maximum number of connections in the pool
	config.MinConns = 5                       // Minimum number of connections in the pool
	config.MaxConnIdleTime = 10 * time.Minute // Maximum idle time before connections are closed

	// Create a connection pool
	DB, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		panic("Unable to connect to database")
	}
}
