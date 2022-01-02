package main

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/desertbit/grumble"
	"nshmadhani.com/notcoin/blockchain/util"
	"nshmadhani.com/notcoin/blockchain/wallet"
	WalletCTX "nshmadhani.com/notcoin/blockchain/wallet"
)

var Manager *WalletManager

func main() {

	InitShell()

	grumble.Main(App)

}

func WalletFromFile(walletFile string) WalletCTX.Wallet {
	if util.FileExists(walletFile) {
		data, _ := ioutil.ReadFile(walletFile)
		return wallet.DeserializeWallet(data)
	} else {
		log.Panic(errors.New("File does not Exist"))
		return WalletCTX.Wallet{}
	}
}
