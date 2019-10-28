package db

import "github.com/shogo82148/txmanager"

func TxExecQuery(dbm txmanager.DB, query string) error {
	err := txmanager.Do(dbm, func(tx txmanager.Tx) error {
		_, err := tx.Query(query)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
