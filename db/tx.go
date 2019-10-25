package db

import "github.com/shogo82148/txmanager"

func TxExecQuery(dbm txmanager.DB, query string) error {
	return txmanager.Do(dbm, func(tx txmanager.Tx) error {
		_, err := tx.Query(query)
		return err
	})
}
