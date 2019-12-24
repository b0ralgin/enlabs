package migration

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20191224082132, Down20191224082132)
}

func Up20191224082132(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE transactions (
                      id varchar(50) NOT NULL,
                      amount integer not null,
                      state varchar(20) not null,
                      source varchar(20) not null,
                      PRIMARY KEY(id))`)
	if err != nil {
		return err
	}
	return nil
}

func Down20191224082132(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE transactions")
	if err != nil {
		return err
	}
	return nil
}
