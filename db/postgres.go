package db

import (
	"database/sql"
	"enlabs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Postgres struct {
	db *sqlx.DB
}

type Keeper interface {
	AddTransaction(t *enlabs.Transaction) error
	GetAmounts() ([]int, error)
}

func NewPostgresClient(dsn string) (*Postgres, error) {
	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "can't create connection to db")
	}
	return &Postgres{conn}, nil
}

func (p *Postgres) AddTransaction(t *enlabs.Transaction) error {
	_, err := p.db.NamedExec(`INSERT INTO transactions (id,amount,state, source) 
			VALUES (:id,:amount,:state, :source)`,
		t)
	if err != nil {
		return errors.Wrap(err, "can't insert transaction")
	}
	return nil
}

func (p *Postgres) GetAmounts() ([]int, error) {
	var trans []int
	update := func(tx *sqlx.Tx) error {
		var balance struct {
			LastID string
			Amount int
		}
		res := tx.QueryRow("SELECT last_id, balance FROM account")
		if err := res.Scan(&balance.LastID, &balance.Amount); err != nil && err != sql.ErrNoRows {
			return errors.Wrap(err, "can't get balance")
		}
		trans = append(trans, balance.Amount)
		err := tx.Select(&trans, "SELECT amount FROM transactions where id > $1", balance.LastID)
		if err != nil {
			return errors.Wrap(err, "can't get transactions")
		}
		return nil
	}
	if err := p.inTx(update); err != nil {
		return nil, errors.Wrap(err, "can't get transactions")
	}
	return trans, nil
}

func (p *Postgres) inTx(fn func(t *sqlx.Tx) error) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "can't create transaction")
	}
	defer tx.Rollback()
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
