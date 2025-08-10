package database

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

func (d *DBWrapper) Set(key, value string) error {
	err := d.DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
	if err != nil {
		return fmt.Errorf("could not insert data: %v", err)
	}
	return nil
}
