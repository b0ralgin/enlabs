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
	err = a.db.InTx(func(tx db.DBTX) error {
		amounts, err := a.db.GetAmounts(tx)
		if err != nil {
			return errors.Wrap(err, "can't get balance")
		}
		res = enlabs.GetBalance(amounts)
		return nil
	})
	return
}

func (a manager) AddTransaction(t *enlabs.Transaction) error {
	if err := a.db.AddTransaction(t); err != nil {
		return errors.Wrap(err, "can't add transaction")
	}
	return nil
}

func (a manager) CorrectBalance() error {
	err := a.db.InTx(func(tx db.DBTX) error {
		if err := a.db.DeleteOddTransactions(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "can't correct balance")
	}
	return nil
}

func (a manager) UpdateBalance() error {
	err := a.db.InTx(func(tx db.DBTX) error {
		amounts, err := a.db.GetAmounts(tx)
		if err != nil {
			return errors.Wrap(err, "can't get balance")
		}
		lastTran := amounts[len(amounts)-1]
		lastTran.Amount = enlabs.GetBalance(amounts)
		if err := a.db.UpdateBalance(tx, lastTran); err != nil {
			return errors.Wrap(err, "can't update balance")
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update balance")
	}
	return nil
}

//Manager account managing interface
type Manager interface {
	GetBalance() (int, error)
	AddTransaction(t *enlabs.Transaction) error
}

type Corrector interface {
	CorrectBalance() error
	UpdateBalance() error
}

//NewAccountManager initialize manager
func NewAccountManager(db db.Keeper) Manager {
	return &manager{db}
}

func NewCorrector(db db.Keeper) Corrector {
	return &manager{db}
}
