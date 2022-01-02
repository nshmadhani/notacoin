package main

import (
	"fmt"
)

func getUrl(address, path string) string {
	return fmt.Sprintf("http://%s%s", address, path)
}

/***********************************************************************************************
*************************************************************************************************
											OLD CODE
***********************************************************************************************
**************************************************************************************************/

/*
	func SendTX(tx tx.Transaction, senderAddress string, amount int, nAddr string) bool {

	transactionBytes := tx.Serialize()

	transactionData := hex.EncodeToString(transactionBytes)

	resp, err := http.PostForm(getUrl(nAddr, "/tx/new"), url.Values{
		"tx":            {transactionData},
		"senderAddress": {senderAddress},
		"amount":        {strconv.Itoa(amount)},
	})
	defer resp.Body.Close()

	if err != nil {
		log.Panic(err)
	}

	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}

	func DownloadBlockchain(address, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(getUrl(address, "/blockchain"))
	if err != nil {
		log.Panic(err)
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}


*/
