package database

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

func GetData(db *badger.DB, key string) string {
	var valueCopy []byte
	err := db.View(func(txn *badger.Txn) error {
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
