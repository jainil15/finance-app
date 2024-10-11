package postgres

import (
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func New(user, password, host, port, dbname string) *sqlx.DB {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname,
	)
	pgxConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		log.Fatalf("Error while parsing database url")
	}
	db, err := sqlx.Connect(
		"pgx",
		pgxConfig.ConnString(),
	)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
