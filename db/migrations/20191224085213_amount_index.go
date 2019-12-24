package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

//nolint:gochecknoinits
func init() {
	goose.AddMigration(Up20191224085213, Down20191224085213)
}

//Up20191224085213 ...
func Up20191224085213(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE INDEX amount_idx ON transactions(id, amount);`)
	if err != nil {
		return err
	}
	return nil
}

//Down20191224085213 ...
func Down20191224085213(tx *sql.Tx) error {
	_, err := tx.Exec("DROP INDEX  amount_idx;")
	if err != nil {
		return err
	}
	return nil
}
