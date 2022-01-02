package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/boltdb/bolt"

	"nshmadhani.com/notcoin/blockchain/block"
	"nshmadhani.com/notcoin/blockchain/network"
	"nshmadhani.com/notcoin/blockchain/tx"
	. "nshmadhani.com/notcoin/blockchain/wallet"
)

const DB_FILE = "wallets.db"
const WALLET_BUCKET = "wallets"

const BLOCKCHAIN_DB = "blockchain.db"

type WalletManager struct {
	WalletDB bolt.DB

	Wallets map[string]*Wallet
	Aliases map[string]string

	Node     *block.Node
	PeerAddr string
}

func (wm *WalletManager) Add(w Wallet) {
	_ = wm.WalletDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(WALLET_BUCKET))
		address := []byte(w.GetAddress())
		err := b.Put(address, w.Serialize())

		fmt.Printf("Added %s=%s\n", w.Alias, w.GetAddress())

		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	wm.UpdateAddress()
}

func (wm *WalletManager) FindWallet(alias string) *Wallet {
	return wm.Wallets[wm.Aliases[alias]]
}

func (wm *WalletManager) Updatechain() {

	// if wm.bc.IsNil() {
	// 	wm.bc = *block.NewBlockchain(db_path)
	// }

	wm.Node.SendGetBlock(wm.PeerAddr)

}

func (wm *WalletManager) UpdateAddress() {
	wallets := make(map[string]*Wallet)
	aliases := make(map[string]string)

	err := wm.WalletDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(WALLET_BUCKET))

		if err != nil {
			log.Panic(err)
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			w := DeserializeWallet(v)

			wallets[w.GetAddress()] = &w
			aliases[w.Alias] = w.GetAddress()
		}

		return nil
	})
	if err != nil {
		log.Panic(err)

	}

	wm.Aliases = aliases
	wm.Wallets = wallets
}
func (wm *WalletManager) SendTX(tx tx.Transaction, nAddr string) {

	transactionBytes := tx.Serialize()

	txComand := network.TX{transactionBytes, ""}

	block.SendCommand(wm.PeerAddr, txComand.Serialize())

}

func LoadManager(peerAddr string) *WalletManager {
	db, err := bolt.Open(DB_FILE, 0666, nil)

	if err != nil {
		log.Panic(err)
		return nil
	}

	wallets := make(map[string]*Wallet)
	aliases := make(map[string]string)

	db_path, _ := filepath.Abs(BLOCKCHAIN_DB)

	node := block.NewNode("9000", db_path, "node.wt")

	go node.Run()
	wm := WalletManager{*db, wallets, aliases, node, peerAddr}

	wm.UpdateAddress()

	wm.Updatechain()

	return &wm
}
