package main

import (
	. "nshmadhani.com/notcoin/blockchain/tx"
)

func CreateNewTransaction(wm WalletManager, fromAlias string, to string, amount int) Transaction {

	fromWallet := wm.FindWallet(fromAlias)

	tx := wm.Node.Chain.NewTransaction(fromWallet.GetAddress(), to, amount, *fromWallet)

	// fmt.Println(tx)

	return *tx
}
