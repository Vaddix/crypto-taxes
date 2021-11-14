package harvest

import (
	"fmt"
)

func Harvest(transactions []chan struct {
	string
	float32
}) float32 {
	//Return the $ capital gains from the given transactions (to the cent)
	fmt.Println("Harvesting tax losses.")
	var capitalGains float32 = 0
	for i := 0; i < len(transactions); i++ {
		transaction := <-transactions[i]
		capitalGains += transaction.float32
	}
	return capitalGains
}
