package db

import (
	"github.com/shogo82148/txmanager"

	"github.com/fatih/color"
)

func TxExecQuery(dbm txmanager.DB, query string, force bool) error {
	err := txmanager.Do(dbm, func(tx txmanager.Tx) error {
		_, err := tx.Query(query)
		if force && err != nil {
			fy := color.New(color.FgYellow)
			fy.Println("Ignore :", err)
			return nil
		}

		if err != nil {
			return err
		}
		return nil
	})
	return err
}
