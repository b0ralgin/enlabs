package account

import (
	"enlabs"
	"enlabs/db"
	"github.com/pkg/errors"
)

type manager struct {
	db db.Keeper
}

func (a manager) GetBalance() (int, error) {
	amounts, err := a.db.GetAmounts()
	if err != nil {
		return 0, errors.Wrap(err, "can't get balance")
	}
	return enlabs.GetBalance(amounts), nil
}

func (a manager) AddTransaction(t *enlabs.Transaction) error {
	if err := a.db.AddTransaction(t); err != nil {
		return errors.Wrap(err, "can't add transaction")
	}
	return nil
}

type Manager interface {
	GetBalance() (int, error)
	AddTransaction(t *enlabs.Transaction) error
}

func NewAccountManager(db db.Keeper) Manager {
	return &manager{db}
}
