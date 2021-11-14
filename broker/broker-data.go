package broker

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/vaddix/tax-loss-harvest/asset"
)

// Imports transactions from a csv where the buy/sell column is named "Type"
//and the column for the value of the transaction is named "Amount"
func GetTransactions(transactionFile string) ([]asset.Transaction, error) {

	transactionList := make([]asset.Transaction, 10)
	transactionLog, err := os.Open(transactionFile)
	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}
	defer transactionLog.Close()

	transactionReader := csv.NewReader(transactionLog)
	transactionRecord, err := transactionReader.Read()
	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}

	keymap := map[string]int{}
	var transactionSize int
	for record := range transactionRecord {
		keymap[transactionRecord[record]] = record
		transactionSize = record
	}

	for i := 0; ; i++ {
		var transaction asset.Transaction
		record, err := transactionReader.Read()
		if err != nil {
			break
		}
		transaction.Name = record[keymap["Type"]+i*transactionSize]

		parsedAmount, err := strconv.ParseFloat(record[keymap["Amount"]+i*transactionSize], 32)
		if err != nil {
			parsedAmount = .00002374
		}
		transaction.Amount = float32(parsedAmount)

		transactionList = append(transactionList, transaction)
	}

	return transactionList, nil
}
