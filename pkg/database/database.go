package database

import (
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var (
	db *DB
)

type DB struct {
	*sqlx.DB
}

func Connect() (*DB, error) {
	if db != nil {
		return db, nil
	}
	// conn, err := sqlx.Connect(sqliteshim.ShimName, "file::memory:?cache=shared&_fk=1")
	conn, err := sqlx.Connect(sqliteshim.ShimName, "file:cake.db?cache=shared&_fk=1")
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		return nil, err
	}
	db = &DB{
		DB: conn,
	}
	f, err := os.ReadFile("sql/schema.sql")
	if err != nil {
		slog.Error("failed to read schema.sql", "err", err)
		return nil, err
	}
	_, err = db.Exec(string(f))
	if err != nil {
		slog.Error("failed to execute schema.sql", "err", err)
		return nil, err
	}
	return db, nil
}
