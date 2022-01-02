package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const targetBits = 12

func toBytes(n int64) []byte {
	return []byte(strconv.FormatInt(n, 10)) //using 16 is like to Hex
}

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewPoW(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {

	return bytes.Join([][]byte{
		pow.block.PrevHash,
		pow.block.HashTXs(),
		toBytes(pow.block.Timestamp),
		toBytes(targetBits),
		toBytes(int64(nonce)),
		toBytes(int64(pow.block.Height)),
	},
		[]byte{})

}

func (pow *ProofOfWork) Validate() bool {

	var cmpInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	cmpHash := sha256.Sum256(data)
	cmpInt.SetBytes(cmpHash[:])

	return cmpInt.Cmp(pow.target) == -1
}

func (pow *ProofOfWork) run() (int, []byte) {

	nonce := 0
	var hashInt big.Int
	var hash [32]byte

	fmt.Printf("Mining a new Block \n")
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		fmt.Printf("\r%x", hash)

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}

	}
	fmt.Printf("\n\n")

	return nonce, hash[:] //[:] What does this dot do? well its an array so we are making it into a slice
}
