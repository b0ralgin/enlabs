package enlabs

type Account struct {
	Balance int
}

func GetBalance(amounts []int) int {
	if len(amounts) == 0 {
		panic("empty amounts array")
	}
	inBalance := amounts[0]
	for _, amount := range amounts[1:] {
		inBalance += amount
		if inBalance < 0 {
			inBalance = 0
		}
	}
	return inBalance
}
