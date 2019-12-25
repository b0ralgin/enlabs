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

type DBTX = *sqlx.Tx

//Keeper interface for repository
type Keeper interface {
	AddTransaction(t *enlabs.Transaction) error
	GetAmounts(tx DBTX) ([]enlabs.Transaction, error)
	DeleteOddTransactions(tx DBTX) error
	UpdateBalance(tx DBTX, t enlabs.Transaction) error
	GetConn() *sql.DB
	InTx(func(tx DBTX) error) error
}

type balance struct {
	LastID int
	Amount int
}

//NewPostgresClient initialize client
func NewPostgresClient(dsn string) (*Postgres, error) {
	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "can't create connection to db")
	}
	return &Postgres{conn}, nil
}

//GetConn get sql.DB struct
func (p *Postgres) GetConn() *sql.DB {
	return p.db.DB
}

func (p *Postgres) GetTx() (*sqlx.Tx, error) {
	return p.db.Beginx()
}

//AddTransaction add transaction to DB
func (p *Postgres) AddTransaction(t *enlabs.Transaction) error {
	_, err := p.db.NamedExec(`INSERT INTO transactions (id,amount,state, source) 
			VALUES (:id,:amount,:state, :source)`,
		t)
	if err != nil {
		return errors.Wrap(err, "can't insert transaction")
	}
	return nil
}

//GetAmounts get amounts of transactions
func (p *Postgres) GetAmounts(tx *sqlx.Tx) ([]enlabs.Transaction, error) {
	var balance balance
	var trans []enlabs.Transaction
	err := tx.Get(&balance, "SELECT last_id as lastId, balance as Amount FROM account ORDER BY last_id DESC LIMIT 1")
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "can't get balance")
	}
	trans = append(trans, enlabs.Transaction{IntID: balance.LastID, Amount: balance.Amount})
	getTranErr := tx.Select(&trans,
		"SELECT internal_id as intID, amount FROM transactions where internal_id > $1 ORDER BY internal_id DESC",
		balance.LastID)
	if getTranErr != nil {
		return nil, errors.Wrap(getTranErr, "can't get transactions")
	}
	return trans, nil
}

func (p *Postgres) DeleteOddTransactions(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE transactions set amount = 0 WHERE internal_id in (SELECT internal_id  FROM transactions 
			WHERE (ROW_NUMBER () OVER (ORDER BY internal_id) %2 = 1 LIMIT 10)`)
	if err != nil {
		return errors.Wrap(err, "can't delete transactions")
	}
	return nil
}

func (p *Postgres) UpdateBalance(tx *sqlx.Tx, t enlabs.Transaction) error {
	_, err := tx.NamedExec(
		"INSERT INTO account (last_id, balance) VALUES (:lastid,  :amount)",
		&balance{LastID: t.IntID, Amount: t.Amount})
	if err != nil {
		return errors.Wrap(err, "can't insert new balance")
	}
	return nil
}

func (p *Postgres) InTx(fn func(t *sqlx.Tx) error) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "can't create transaction")
	}
	defer func() {
		_ = tx.Rollback()
	}()
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
