package asset

type Transaction struct {
	Name   string
	Type   string
	Amount float32
	Basis  float32
}

type Holding struct {
	Name    string
	History []Transaction
}
