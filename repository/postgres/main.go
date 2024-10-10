package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	db, err := sqlx.Connect(
		"postgres",
		connectionString,
	)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
