package database

import (
	"os"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v4"
)

type DBWrapper struct {
	DB *badger.DB
}

func OpenDB() (*DBWrapper, error) {
	var dbPath string
	
	// Use /tmp for development, ~/.timetrack for production
	if os.Getenv("TIMETRACK_ENV") == "development" {
		dbPath = "/tmp/badger"
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dbPath = filepath.Join(homeDir, ".timetrack", "data")
		
		// Create directory if it doesn't exist
		if err := os.MkdirAll(dbPath, 0755); err != nil {
			return nil, err
		}
	}
	
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
