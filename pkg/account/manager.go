package account

import (
	"enlabs"
	"enlabs/db"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type manager struct {
	db  db.Keeper
	log *logrus.Entry
}

func (a manager) GetBalance() (res int, err error) {
	balance, err := a.db.GetBalance()
	if err != nil {
		a.log.WithError(err).Error("can't get balance")
		return 0, err
	}
	return balance, nil
}

func (a manager) AddTransaction(t *enlabs.Transaction) error {
	err := a.db.AddTransaction(t)
	if err != nil {
		a.log.WithError(err).Error("can't add transaction")
		return err
	}
	return nil
}

func (a manager) CorrectBalance() error {
	l := logrus.WithField("method", "correct balance")
	tx, err := a.db.BeginTx()
	if err != nil {
		l.WithError(err).Error("can't create tx")
		return err
	}
	defer tx.Rollback()
	balance, err := tx.GetBalance()
	if err != nil {
		l.WithError(err).Error("can't get balance")
		return err
	}
	trans, err := tx.GetLastN(10)
	if err != nil {
		l.WithError(err).Error("can't get transactions")
	}
	for i, t := range trans {
		if i%2 == 0 {
			continue
		}
		if err := tx.DeleteTransaction(t); err != nil {
			l.WithError(err).Errorf("can't delete transaction %s", t.ID)
			return err
		}
		t.Amount = -t.Amount
		balance = t.CalcBalance(balance)
	}
	if err := tx.UpdateBalance(balance); err != nil {
		l.WithError(err).Error("can't update balance")
		return err
	}
	return errors.Wrap(tx.Commit(), "can't commit balance correction")
}

//Manager account managing interface
type Manager interface {
	GetBalance() (int, error)
	AddTransaction(t *enlabs.Transaction) error
}

//Corrector interface for manager
type Corrector interface {
	CorrectBalance() error
}

//NewAccountManager initialize manager
func NewAccountManager(db db.Keeper, log *logrus.Entry) Manager {
	return &manager{db, log.WithField("service", "manager")}
}

//NewCorrector initialize corrector manager
func NewCorrector(db db.Keeper, log *logrus.Entry) Corrector {
	return &manager{db, log.WithField("service", "corrector")}
}
