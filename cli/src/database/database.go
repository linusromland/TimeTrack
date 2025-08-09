package database

import (
	badger "github.com/dgraph-io/badger/v4"
)

type DBWrapper struct {
	DB *badger.DB
}

func OpenDB() (*DBWrapper, error) {
	dbPath := "/tmp/badger"
	opts := badger.DefaultOptions(dbPath).WithInMemory(false)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &DBWrapper{DB: db}, nil
}

func (d *DBWrapper) Close() error {
	return d.DB.Close()
}
