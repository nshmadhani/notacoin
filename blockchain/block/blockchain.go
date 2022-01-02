package block

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
	. "nshmadhani.com/notcoin/blockchain/tx"
	"nshmadhani.com/notcoin/blockchain/util"
	. "nshmadhani.com/notcoin/blockchain/wallet"
)

// const dbfile = "blockchain.DB"
const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte
	DB  *bolt.DB
}

/******************************************************************************************************************
													Block
******************************************************************************************************************/

func (blockchain *Blockchain) MineBlock(transactions []*Transaction) {
	block := NewBlock(transactions, blockchain.tip, blockchain.GetBestHeight()+1)
	blockchain.AddBlock(block)
}

func (blockchain *Blockchain) AddBlock(block *Block) {

	if bytes.Compare(block.PrevHash, blockchain.tip) != 0 {
		fmt.Printf("Did not add Block=%s on tip=%s\n", hex.EncodeToString(block.Hash), hex.EncodeToString(blockchain.tip))
		return
	}

	err := blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		_ = b.Put(block.Hash, block.Serialize())

		err := b.Put([]byte("l"), block.Hash)

		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		blockchain.tip = block.Hash

		return err
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	utxoSet := UTXO{blockchain}
	utxoSet.Reindex()

}

func (blockchain Blockchain) GetBlockchainHashs() [][]byte {

	bci := blockchain.Iterator()
	var hashs [][]byte
	for {
		block := bci.Next()
		hashs = append(hashs, block.Hash)
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return hashs
}

func (blockchain Blockchain) GetBlock(blockHash []byte) Block {
	var block Block
	_ = blockchain.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(blockHash)

		if len(encodedBlock) != 0 {
			block = *DeserializeBlock(encodedBlock)
		}
		return nil
	})
	return block
}

/******************************************************************************************************************
							Blockchain Creation
******************************************************************************************************************/

func NewBlockchain(DB_FILE string) *Blockchain {

	if !util.FileExists(DB_FILE) {
		return CreateBlockchain(DB_FILE)
	}

	db, err := bolt.Open(DB_FILE, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte

	err = db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(blocksBucket))

		if bucket == nil {
			return errors.New("EMPTY")
		} else {
			tip = bucket.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
		defer db.Close()
		return &Blockchain{}
	}

	bc := &Blockchain{tip, db}
	utxoSet := UTXO{bc}
	utxoSet.Reindex()

	return bc
}

func CreateBlockchain(dbfile string) *Blockchain {

	db, err := bolt.Open(dbfile, 0600, nil)

	if err != nil {

		log.Fatal(err)
	}

	var tip []byte

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return err
		}
		cbTx := NewCoinBaseTx("M28juHW58uH7K18bZPES8NzyjaFjSAPeRP", "Video Killed the Radio Star", 0)

		fmt.Println(cbTx)

		gen := GenesisBlock(cbTx)

		tip = gen.Hash

		bucket.Put(gen.Hash, gen.Serialize())
		bucket.Put([]byte("l"), gen.Hash)

		return nil
	})
	if err != nil {
		log.Panic(err)
		defer db.Close()
		return nil
	}
	bc := &Blockchain{tip, db}
	utxoSet := UTXO{bc}
	utxoSet.Reindex()

	return bc

}

/******************************************************************************************************************
							Blockchain Transaction
******************************************************************************************************************/

func (bc Blockchain) FindUTXO() map[string][]*TxOutput {

	utxos := make(map[string][]*TxOutput)
	utxosMap, utxs, _ := bc.FindUTXOs([]byte{}, 0)
	for txId, tx := range utxs {
		for _, idx := range utxosMap[txId] {
			utxos[txId] = append(utxos[txId], &tx.Vout[idx])
		}
	}
	return utxos
}

func (bc Blockchain) FindUTXOs(pubKeyHash []byte, amount int) (map[string][]int, map[string]*Transaction, int) {

	spentTXOs := make(map[string][]int)
	unspentTXs := make(map[string]*Transaction)
	unspentTXOs := make(map[string][]int)

	bci := bc.Iterator()
	acc := 0

Blocks:
	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID)
			if !tx.IsCoinbase() {
				for _, txin := range tx.Vin {
					if len(pubKeyHash) == 0 || txin.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(txin.TxId) // the txID to which the output belongs to
						spentTXOs[inTxID] = append(spentTXOs[inTxID], txin.Vout)
					}
				}
			}

		Outputs:
			for outIdx, out := range tx.Vout {

				if len(pubKeyHash) != 0 && !out.IsLockedWithKey(pubKeyHash) {
					continue Outputs
				}

				if spentTXOs[txId] != nil {
					for _, vout := range spentTXOs[txId] {
						if vout == outIdx {

							continue Outputs
						}
					}
				}

				unspentTXOs[txId] = append(unspentTXOs[txId], outIdx)
				acc += out.Value

				if unspentTXs[txId] == nil {
					unspentTXs[txId] = tx
				}

				if acc > amount && amount > 0 {
					fmt.Printf("Broke acc=%d,amount=%d\n", acc, amount)
					break Blocks
				}

			}

		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTXOs, unspentTXs, acc
}

func (bc Blockchain) NewTransaction(from string, to string, amount int, w Wallet) *Transaction {

	pubKeyHash := util.HashPubKey(w.Publickey)

	utxoSet := UTXO{Blockchain: &bc}

	utxos, acc := utxoSet.FindUTXO(pubKeyHash, amount)

	if acc < amount {
		fmt.Printf("acc=%d,amt=%d\n", acc, amount)

		log.Panic("Not Enough Balance")
		return nil
	}

	var inputs []TxInput
	var outputs []TxOutput

	for txId := range utxos {
		for _, out := range utxos[txId] {
			a, _ := hex.DecodeString(txId)
			inputs = append(inputs, TxInput{out.Index, a, nil, w.Publickey})
		}
	}

	outputs = append(outputs, NewTxOuput(amount, to, 0))

	if acc > amount {
		outputs = append(outputs, NewTxOuput(acc-amount, from, 1))
	}

	tx := Transaction{nil, inputs, outputs}

	w.Sign(&tx, utxos)
	tx.SetId()

	fmt.Println(utxos)

	fmt.Println("Verification is", bc.VerifyTx(&tx))

	return &tx
}

func (bc *Blockchain) VerifyTx(tx *Transaction) bool {

	//TO Verify a transaction take a look at its

	utxos := make(map[string][]*TxOutput)

	utxoSet := NewUTXO(bc)

	for _, vin := range tx.Vin {

		txId := hex.EncodeToString(vin.TxId)

		out := utxoSet.FindTx(vin.TxId, vin.Vout)
		if out == nil {
			return false
		}
		utxos[txId] = append(utxos[txId], out)
	}
	fmt.Println("Verification UTXO used", utxos)

	return tx.Verify(utxos)

}

/*  ===========================================================================================================
								MISC
===========================================================================================================
*/

func (bc *Blockchain) Close() {
	bc.DB.Close()
}

func (bc Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

func (bc *Blockchain) IsNil() bool {
	return len(bc.tip) == 0
}

/*  ===========================================================================================================
								Blockchain Iterator
===========================================================================================================
*/

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.DB}
}

type BlockchainIterator struct {
	currentHash []byte
	DB          *bolt.DB
}

func (bci *BlockchainIterator) Next() *Block {

	var block *Block

	err := bci.DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(bci.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
		return nil
	}

	bci.currentHash = block.PrevHash

	return block

}
