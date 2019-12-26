package db

import (
	"database/sql"
	"enlabs"

	"github.com/jmoiron/sqlx"
	//sql driver import
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

//Postgres client of DB
type Postgres struct {
	db *sqlx.DB
}

//Keeper interface for repository
type Keeper interface {
	AddTransaction(t *enlabs.Transaction) error
	GetBalance() (int, error)
	BeginTx() (TXer, error)
}

//NewPostgresClient initialize client
func NewPostgresClient(dsn string) (*Postgres, error) {
	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "can't create connection to db")
	}
	return &Postgres{conn}, nil
}

func (p *Postgres) BeginTx() (TXer, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "can't create tx")
	}
	return &TX{tx}, nil
}

func (p *Postgres) GetConn() *sql.DB {
	return p.db.DB
}

//AddTransaction add transaction to DB
func (p *Postgres) AddTransaction(t *enlabs.Transaction) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}
	trx := TX{tx}
	defer tx.Rollback()
	balance, getBalErr := trx.GetBalance()
	if getBalErr != nil {
		return getBalErr
	}
	_, addTranErr := tx.NamedExec(`INSERT INTO transactions (id,amount,state, source) 
			VALUES (:id,:amount,:state, :source)`,
		t)
	if addTranErr != nil {
		return errors.Wrap(addTranErr, "can't insert transaction")
	}
	if err := trx.UpdateBalance(t.CalcBalance(balance)); err != nil {
		return err
	}
	return tx.Commit()
}

//GetAmounts get amounts of transactions
func (p *Postgres) GetBalance() (int, error) {
	tx, err := p.BeginTx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	return tx.GetBalance()
}
