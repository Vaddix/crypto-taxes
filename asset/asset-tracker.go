package asset

type Transaction struct {
	Name   string
	Amount float32
}

type Holding struct {
	Name    string
	History []Transaction
}
