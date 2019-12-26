package account

import (
	"enlabs"
	"enlabs/db"

	"github.com/pkg/errors"
)

type manager struct {
	db db.Keeper
}

func (a manager) GetBalance() (res int, err error) {
	balance, err := a.db.GetBalance()
	if err != nil {
		return 0, errors.Wrap(err, "can't get balance")
	}
	return balance, nil
}

func (a manager) AddTransaction(t *enlabs.Transaction) error {
	if err := a.db.AddTransaction(t); err != nil {
		return errors.Wrap(err, "can't add transaction")
	}
	return nil
}

func (a manager) CorrectBalance() error {
	tx, err := a.db.BeginTx()
	if err != nil {
		return errors.Wrap(err, "can't create tx for balance correction")
	}
	defer tx.Rollback()
	balance, err := tx.GetBalance()
	if err != nil {
		return errors.Wrap(err, "can't get balance")
	}
	trans, err := tx.GetLastN(10)
	if err != nil {
		return errors.Wrap(err, "can't get transactions")
	}
	for i, t := range trans {
		if i%2 == 0 {
			continue
		}
		if err := tx.DeleteTransaction(t); err != nil {
			return errors.Wrapf(err, "can't delete transaction %s", t.ID)
		}
		t.Amount = -t.Amount
		balance = t.CalcBalance(balance)
	}
	if err := tx.UpdateBalance(balance); err != nil {
		return errors.Wrap(err, "can't update balance")
	}
	return errors.Wrap(tx.Commit(), "can't commit balance correction")
}

//Manager account managing interface
type Manager interface {
	GetBalance() (int, error)
	AddTransaction(t *enlabs.Transaction) error
}

type Corrector interface {
	CorrectBalance() error
}

//NewAccountManager initialize manager
func NewAccountManager(db db.Keeper) Manager {
	return &manager{db}
}

func NewCorrector(db db.Keeper) Corrector {
	return &manager{db}
}
