package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

//nolint:gochecknoinits
func init() {
	goose.AddMigration(Up20191224082904, Down20191224082904)
}

//Up20191224082904 ...
func Up20191224082904(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE account (
    updated_at timestamp not null default current_timestamp ,
    balance integer not null
    )`)
	if err != nil {
		return err
	}
	return nil
}

//Down20191224082904 ...
func Down20191224082904(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE account")
	if err != nil {
		return err
	}
	return nil
}
