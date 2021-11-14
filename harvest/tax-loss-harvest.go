package harvest

import (
	"fmt"

	"github.com/vaddix/tax-loss-harvest/asset"
)

func Harvest(transactions []asset.Transaction) float32 {
	//Return the $ capital gains from the given transactions (to the cent)
	fmt.Println("Harvesting tax losses.")
	var capitalGains float32 = 0
	for i := 0; i < len(transactions); i++ {
		capitalGains += transactions[i].Amount
	}
	return capitalGains
}
