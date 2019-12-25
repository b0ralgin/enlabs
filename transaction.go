package enlabs

//State state of transaction
type State string

//Win increase balance
const Win = State("win")

//Lost decrease balance
const Lost = State("lost")

//Source source of transaction
type Source string

//ServerSource from server
const ServerSource = "server"

//ClientSource from client
const ClientSource = "client"

//PaymentSource from payment
const PaymentSource = "payment"

//Transaction model
type Transaction struct {
	IntID  int
	ID     string
	Amount int
	State  State
	Source Source
}

//NewTransaction create new transaction
func NewTransaction(id string, amount int, state State, source Source) *Transaction {
	if state == Lost {
		amount = -amount
	}
	return &Transaction{
		ID:     id,
		Amount: amount,
		State:  state,
		Source: source,
	}
}
