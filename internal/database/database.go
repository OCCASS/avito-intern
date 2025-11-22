package database

import (
	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	Conn *sqlx.DB
}

func MustConnect(conn string) *Database {
	db, err := sqlx.Connect("pgx", conn)
	if err != nil {
		panic(err)
	}

	return &Database{
		Conn: db,
	}
}
