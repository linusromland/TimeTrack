package database

import (
	badger "github.com/dgraph-io/badger/v4"
)

func OpenDB() (*badger.DB, error) {
	dbPath := "/tmp/badger"
	opts := badger.DefaultOptions(dbPath).WithInMemory(true)
	opts.InMemory = false
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDB(db *badger.DB) error {
	return db.Close()
}
