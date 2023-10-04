package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/jadwahab/go-prescription-airdrop/helpers"
	"github.com/mrz1836/go-whatsonchain"
)

func main() {
	filename := "perscList.json"
	perscList, err := readFile(filename)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	client := whatsonchain.NewClient(whatsonchain.NetworkMain, nil, nil)

	var perscListSuccess []Prescription

	for _, persc := range perscList {
		txhex, err := helpers.PrescriptionAirdrop(
			persc.OwnerAddress,
			"e027531de98a7ee6830a45d78e8f2ae873640342fab9190365fb0a27d57ac69a", // inscUTXO
			"L2fgiaz4bpCdvbHvvTgeG4gHC15QFdbwMbwsoyQNTwGVJK4Uido6",             // inscWIF
			"KwPVoki5qyur6JwotJkNXsEEENf46VmtHbMnLESatNoeVhZ1NEGp",             // fundingWIF
			client,
		)
		fmt.Println(txhex)

		if err == nil {
			persc.AirDropTx = &txhex
			perscListSuccess = append(perscListSuccess, persc)
			fmt.Println("Success airdropping:")
			fmt.Println(persc)

			delay(100)
		} else {
			fmt.Println(err)
		}
	}

	// Filter out the successful ones from the original array
	for i, persc := range perscList {
		for _, successPersc := range perscListSuccess {
			if persc == successPersc {
				perscList = append(perscList[:i], perscList[i+1:]...)
				break
			}
		}
	}

	fmt.Println("Successful List:")
	fmt.Println(perscListSuccess)
	fmt.Println("Remaining List:")
	fmt.Println(perscList)

	// Write perscListSuccess and perscList to separate files
	if err := writeToFile("perscListSuccess", perscListSuccess); err != nil {
		fmt.Println("Error writing to success file:", err)
	}

	if err := writeToFile("perscList", perscList); err != nil {
		fmt.Println("Error writing to remaining file:", err)
	}
}

type Prescription struct {
	PropsNo      int     `json:"propsNo"`
	TxId         string  `json:"txid"`
	Location     string  `json:"location"`
	Origin       string  `json:"origin"`
	Seller       *string `json:"seller,omitempty"`
	OwnerAddress string  `json:"ownerAddress"`
	Paymail      string  `json:"paymail"`
	AirDropTx    *string `json:"airdropTx,omitempty"`
}

func readFile(filename string) ([]Prescription, error) {
	// Read the JSON file
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a slice of Prescription structs
	var prescriptions []Prescription
	if err := json.Unmarshal(fileContent, &prescriptions); err != nil {
		return nil, err
	}

	return prescriptions, nil
}

func writeToFile(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, jsonData, 0644)
}

func delay(milliseconds time.Duration) {
	time.Sleep(milliseconds * time.Millisecond)
}
