package block

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	. "nshmadhani.com/notcoin/blockchain/tx"
)

type Block struct {
	Timestamp    int64
	Transactions []*Transaction
	PrevHash     []byte
	Hash         []byte
	Nonce        int
	Height       int
}

func (b *Block) SetHash() {

	pow := NewPoW(b)
	nonce, h := pow.run()
	b.Nonce = nonce
	b.Hash = h[:]
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func (b *Block) HashTXs() []byte {
	tree := CreateMerkleTree(b.Transactions)
	return tree.Root.Data
}

func (b Block) String() string {
	printStr := ""
	printStr += fmt.Sprintf("Nonce is %d\n", b.Nonce)
	printStr += fmt.Sprintf("Timestamp is %d\n", b.Timestamp)
	printStr += fmt.Sprintf("PrevHash is %s\n", hex.EncodeToString(b.PrevHash))
	printStr += fmt.Sprintf("Hash is %s\n\nTransactions: \n\n", hex.EncodeToString(b.Hash))

	for _, tx := range b.Transactions {
		printStr += tx.String() + "\n\n"
	}
	return printStr

}

func (b Block) Less(that *Block) bool {
	if b.Height == that.Height {
		return b.Timestamp < that.Timestamp
	}
	return b.Height < that.Height
}

func DeserializeBlock(blockBytes []byte) *Block {

	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
		return &block
	}
	return &block
}

func NewBlock(txs []*Transaction, prevHash []byte, height int) *Block {
	b := &Block{time.Now().Unix(), txs, prevHash, []byte{}, 0, height}
	b.SetHash()
	return b
}

func GenesisBlock(coinbase *Transaction) *Block {
	b := &Block{1438269960, []*Transaction{coinbase}, []byte{}, []byte{}, 0, 1}
	b.SetHash()
	return b
}
