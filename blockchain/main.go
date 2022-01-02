package main

import (
	"fmt"

	"nshmadhani.com/notcoin/blockchain/tx"
)

func main() {

	aTx := tx.NewCoinBaseTx("M7N5kd1jZt8NM9z7oQzRbhrHLH7CPCSZ7Q", "sdasd", 2)

	aTx.SetId()

	fmt.Println(aTx.String())

}
