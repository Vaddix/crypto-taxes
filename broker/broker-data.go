package broker

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/vaddix/tax-loss-harvest/asset"
)

// Imports transactions from a csv where the buy/sell column is named "Type"
//and the column for the value of the transaction is named "Amount"
func GetTransactions(eventFile string) ([]asset.Transaction, error) {

	transactionList := make([]asset.Transaction, 10)
	eventSize := 10

	eventReader, keymap, _ := ReadEvents(eventFile, eventSize)

	for i := 0; ; i++ {
		record, err := eventReader.Read()
		if err != nil {
			break
		}

		transactionList = append(transactionList, ParseTransaction(record, keymap)...)
	}

	return transactionList, nil
}

func ParseTransaction(record []string, keymap map[string]int) []asset.Transaction {

	var transaction asset.Transaction

	transaction.Name = record[keymap["Asset"]]
	parsedQuantity, _ := strconv.ParseFloat(record[keymap["Quantity Transacted"]], 32)
	transaction.Amount = float32(parsedQuantity)
	parsedQuantity, _ = strconv.ParseFloat(record[keymap["Spot Price"]], 32)
	transaction.Basis = float32(parsedQuantity)
	eventType := record[keymap["Transaction Type"]]
	if eventType == "Buy" || eventType == "Receive" || eventType == "Rewards Income" {
		transaction.Type = "Buy"
	} else if eventType == "Sell" {
		transaction.Type = "Sell"
	} else if eventType == "Convert" {
		return HandleConversionEvent(record, keymap)
	}

	return []asset.Transaction{transaction}
}

func HandleConversionEvent(record []string, keymap map[string]int) []asset.Transaction {
	eventDetails := record[keymap["Notes"]]

	noteList := strings.Split(eventDetails, " ")
	soldAmount := noteList[1]
	soldAsset := noteList[2]
	boughtAmount := noteList[4]
	boughtAsset := noteList[5]

	buyRecord := make([]string, 0, len(record))
	sellRecord := make([]string, 0, len(record))

	copy(buyRecord, record)
	copy(sellRecord, record)

	buyRecord[keymap["Transaction Type"]] = "Buy"
	buyRecord[keymap["Quantity Transacted"]] = boughtAmount
	buyRecord[keymap["Asset"]] = boughtAsset
	sellRecord[keymap["Transaction Type"]] = "Sell"
	sellRecord[keymap["Quantity Transacted"]] = soldAmount
	sellRecord[keymap["Asset"]] = soldAsset

	return append(ParseTransaction(buyRecord, keymap), ParseTransaction(sellRecord, keymap)...)
}

func ReadEvents(eventFile string, eventSize int) (*csv.Reader, map[string]int, error) {

	transactionLog, err := os.Open(eventFile)
	if err != nil {
		fmt.Println("Error:", err)
		transactionLog.Close()
		return nil, nil, errors.New("Failed to open file with given name: " + eventFile)
	}
	defer transactionLog.Close()

	transactionReader := csv.NewReader(transactionLog)
	transactionReader.FieldsPerRecord = eventSize

	for i := 0; true; i++ {
		transactionRecord, err := transactionReader.Read()
		if err != nil {
			fmt.Println("Error:", err)
			return nil, nil, errors.New("Error occurred while reading transaction record " + strconv.FormatInt(int64(i), 10) + ".")
		}
		fmt.Println("Found record: ", transactionRecord)
		if len(transactionRecord) > 3 {
			return transactionReader, GetKeyMap(transactionRecord), nil
		}
		if i > 10 {
			break
		}
	}
	transactionLog.Close()
	return transactionReader, make(map[string]int), nil

}

func GetKeyMap(transactionRecord []string) map[string]int {
	keymap := map[string]int{}

	for recordIndex := range transactionRecord {
		keymap[transactionRecord[recordIndex]] = recordIndex
	}

	return keymap
}
