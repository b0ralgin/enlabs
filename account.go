package enlabs

//Account account model
type Account struct {
	Balance int
}

//GetBalance count balance
func GetBalance(trans []Transaction) int {
	if len(trans) == 0 {
		panic("empty amounts array")
	}
	newBalance := trans[0].Amount
	for _, tran := range trans {
		newBalance += tran.Amount
		if newBalance < 0 {
			newBalance = 0
		}
	}
	return newBalance
}
