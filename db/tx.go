package db

import (
	"database/sql"
	"enlabs"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//TX ...
type TX struct {
	*sqlx.Tx
}

//TXer interface for
type TXer interface {
	GetBalance() (int, error)
	GetLastN(n int) ([]enlabs.Transaction, error)
	DeleteTransaction(t enlabs.Transaction) error
	UpdateBalance(balance int) error
	Rollback()
	Commit() error
}

//Rollback rollback tx
func (tx TX) Rollback() {
	err := tx.Tx.Rollback()
	logrus.Error(err)
}

//GetBalance get balance of account
func (tx TX) GetBalance() (int, error) {
	var balance int
	if err := tx.Get(&balance, "SELECT balance from account LIMIT 1 FOR UPDATE"); err != nil && err != sql.ErrNoRows {
		return 0, errors.Wrap(err, "can't get balance")
	}
	return balance, nil
}

//UpdateBalance update balance
func (tx TX) UpdateBalance(balance int) error {
	_, updateBalErr := tx.NamedExec("UPDATE account SET balance = :amount",
		map[string]interface{}{
			"amount": balance,
		})
	if updateBalErr != nil {
		return errors.Wrap(updateBalErr, "can't update balance")
	}
	return nil
}

//GetLastN get last n records of transactions
func (tx TX) GetLastN(n int) (trans []enlabs.Transaction, err error) {
	err = tx.Select(&trans, "SELECT id, amount FROM transactions order by created_at DESC LIMIT $1", n)
	if err != nil {
		return nil, errors.Wrap(err, "can't get odd transactions")
	}
	return
}

//DeleteTransaction delete transaction
func (tx TX) DeleteTransaction(tran enlabs.Transaction) error {
	if _, err := tx.NamedExec("UPDATE transactions set amount = 0 where id = :id", tran); err != nil {
		return errors.Wrapf(err, "can't delete transactions ")
	}
	return nil
}
