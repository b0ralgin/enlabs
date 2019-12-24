package api

import (
	"enlabs"
	"errors"
	"strconv"
	"strings"
)

type addTransactionRequest struct {
	ID     string       `json:"transactionId"`
	Amount string       `json:"amount"`
	State  enlabs.State `json:"state"`
}

func (r addTransactionRequest) MapToTransaction(source enlabs.Source) (*enlabs.Transaction, error) {
	amount, err := parseAmount(r.Amount)
	if err != nil {
		return nil, err
	}
	return enlabs.NewTransaction(r.ID, int(amount), r.State, source), nil
}

func parseAmount(m string) (int64, error) {
	var cents string
	switch len(strings.Split(m, ".")) {
	case 1:
		cents = m + "00"
	case 2:
		cents = strings.ReplaceAll(m, ".", "")
	default:
		return 0, errors.New("wrong type of amount")
	}
	return strconv.ParseInt(cents, 10, 64)
}
