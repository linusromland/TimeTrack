package database

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

func (d *DBWrapper) Get(key string) string {
	var valueCopy []byte
	err := d.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return fmt.Errorf("could not get data: %v", err)
		}
		valueCopy, err = item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("could not get data: %v", err)
		}
		return nil
	})
	if err != nil {
		return ""
	}
	return string(valueCopy)
}
