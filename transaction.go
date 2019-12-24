package enlabs

type State string

const Win = State("win")
const Lost = State("lost")

type Source string

const ServerSource = "server"
const ClientSource = "client"
const PaymentSource = "payment"

type Transaction struct {
	ID     string
	Amount int
	State  State
	Source Source
}

func NewTransaction(id string, amount int, state State, source Source) *Transaction {
	return &Transaction{
		ID:     id,
		Amount: amount,
		State:  state,
		Source: source,
	}
}
