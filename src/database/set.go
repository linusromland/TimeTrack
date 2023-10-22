package database

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

func SetData(db *badger.DB, key, value string) error {
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
	if err != nil {
		return fmt.Errorf("could not insert data: %v", err)
	}
	return nil
}
