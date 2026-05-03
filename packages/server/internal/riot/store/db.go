package store

import platformdb "github.com/galchammat/kadeem/internal/platform/database"

type Store struct {
	db *platformdb.DB
}

func New(db *platformdb.DB) *Store {
	return &Store{db: db}
}
