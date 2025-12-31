package store

import "database/sql"

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(username string) error {
	// TODO: Generate ULID
	// TODO: INSERT INTO users
	return nil
}
