package postgres

import platformdb "github.com/galchammat/kadeem/internal/platform/database"

type DB struct {
	db *platformdb.DB
}

func New(db *platformdb.DB) *DB {
	return &DB{db: db}
}
