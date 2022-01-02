package tx

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

func (tx *Transaction) SetId() {
	tx.ID = tx.Hash()
}

func (tx *Transaction) Bytes() []byte {

	b := tx.ID

	for _, vin := range tx.Vin {
		b = append(b, vin.Bytes()...)
	}
	for _, vout := range tx.Vout {
		b = append(b, vout.Bytes()...)
	}
	return b

}

func (tx Transaction) Hash() []byte {
	var hash [32]byte

	// We used to create copies before but now, we dont instead hash the same object, as we will set the ID anyway
	// txCopy := *tx
	// txCopy.ID = []byte{}
	tx.ID = []byte{}
	hash = sha256.Sum256(tx.Bytes())

	return hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxId) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) Trim() *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Vin {
		inputs = append(inputs, TxInput{in.Vout, in.TxId, nil, nil})
	}
	for _, out := range tx.Vout {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash, out.Index})
	}
	return &Transaction{tx.ID, inputs, outputs}
}

func (tx *Transaction) Verify(prevTXs map[string][]*TxOutput) bool {

	if tx.IsCoinbase() {
		return true
	}
	trimmedCopy := tx.Trim()
	curve := elliptic.P256()

	for iDx, in := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(in.TxId)]
		trimmedCopy.Vin[iDx].Signature = nil

		for _, out := range prevTx {
			if out.Index == in.Vout {
				trimmedCopy.Vin[iDx].PubKey = out.PubKeyHash
			}
		}

		trimmedCopy.SetId()
		trimmedCopy.Vin[iDx].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}

		fmt.Println(trimmedCopy.String())

		if ecdsa.Verify(&rawPubKey, trimmedCopy.ID, &r, &s) == false {

			fmt.Printf("Verification Failed\n %s", in.String())

			return false
		}

	}
	return true
}

func (b *Transaction) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func (tx Transaction) String() string {
	str := "ID = " + hex.EncodeToString(tx.ID) + "\n"
	for idx, vin := range tx.Vin {
		str += fmt.Sprintf("Input %d\n", idx) + vin.String() + "\n"
	}
	for idx, vout := range tx.Vout {
		str += fmt.Sprintf("Output %d\n", idx) + vout.String() + "\n"
	}
	return str
}

func DeserializeTX(blockBytes []byte) *Transaction {

	var block Transaction
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}

func NewCoinBaseTx(to, data string, index int) *Transaction {

	txIn := TxInput{-1, []byte{}, nil, []byte(data)}

	tx := &Transaction{nil, []TxInput{txIn}, []TxOutput{NewTxOuput(subsidy, to, index)}}

	tx.SetId()

	return tx

}
