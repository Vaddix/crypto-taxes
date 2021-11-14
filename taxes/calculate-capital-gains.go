package taxes

import (
	"fmt"

	"github.com/vaddix/tax-loss-harvest/asset"
	"github.com/vaddix/tax-loss-harvest/harvest"
)

func CalculateCapitalGains(holdings []asset.Holding) float32 {
	var netCapitalGains float32 = 0
	for i := 0; i < len(holdings); i++ {
		gain := harvest.Harvest(holdings[i].History)
		netCapitalGains += gain

		fmt.Println("Capital gain of", gain, "for holding:", holdings[i].Name)
	}
	return netCapitalGains
}
