package db

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func InitDb(conn string) (*sqlx.DB, error) {
	return sqlx.Connect("pgx", conn)
}

func CloseDb(db *sqlx.DB) error {
	return db.Close()
}

func InitialCleanup(db *sqlx.DB) error {
	_, err := db.Queryx("DELETE FROM slugs")
	return err
}
