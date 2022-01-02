package block

import (
	"encoding/hex"
	"log"

	"go.etcd.io/bbolt"
	TransactionCTX "nshmadhani.com/notcoin/blockchain/tx"
)

const utxobucket = "utxo"

type UTXO struct {
	Blockchain *Blockchain
}

func (UTXOSet *UTXO) Reindex() {

	err := UTXOSet.Blockchain.DB.Update(func(tx *bbolt.Tx) error {
		err := tx.DeleteBucket([]byte(utxobucket))
		_, err = tx.CreateBucket([]byte(utxobucket))
		return err

	})

	utxos := UTXOSet.Blockchain.FindUTXO()

	err = UTXOSet.Blockchain.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(utxobucket))
		var err error
		for txID, outs := range utxos {
			key, _ := hex.DecodeString(txID)

			err = b.Put(key, TransactionCTX.SerializeOutputs(outs))

		}
		return err
	})

	if err != nil {
		log.Panic(err)
	}
}

func (UTXOSet *UTXO) FindUTXO(pubKeyHash []byte, amount int) (map[string][]*TransactionCTX.TxOutput, int) {

	outs := make(map[string][]*TransactionCTX.TxOutput)
	acc := 0
	err := UTXOSet.Blockchain.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(utxobucket))
		c := b.Cursor()

	Outputs:
		for k, v := c.First(); k != nil; k, v = c.Next() {

			txId := hex.EncodeToString(k)
			outputs := TransactionCTX.DeserializeTxOutput(v)

			for _, out := range outputs {
				if out.IsLockedWithKey(pubKeyHash) && (amount <= 0 || acc < amount) {
					acc += out.Value
					outs[txId] = append(outs[txId], out)
				} else if acc > amount {
					break Outputs
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return outs, acc
}

func (UTXOSet *UTXO) Update(block Block) {

	updatedOutputs := make(map[string][]*TransactionCTX.TxOutput)
	err := UTXOSet.Blockchain.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(utxobucket))

		for _, tx := range block.Transactions {

			if !tx.IsCoinbase() {

				for _, vin := range tx.Vin {

					outputTxId := hex.EncodeToString(vin.TxId)

					outputIdx := vin.Vout

					if len(updatedOutputs[outputTxId]) == 0 {
						//load from database
						byteOutputs := b.Get(vin.TxId)
						updatedOutputs[outputTxId] = append(updatedOutputs[outputTxId],
							TransactionCTX.DeserializeTxOutput(byteOutputs)...)
					}

					// Delete that particular Index

					updatedOutputs[outputTxId] = deleteVout(updatedOutputs[outputTxId], outputIdx)

					if len(updatedOutputs[outputTxId]) == 0 {
						b.Delete(vin.TxId)
					}

				}
			}
			for _, out := range tx.Vout {
				updatedOutputs[hex.EncodeToString(tx.ID)] = append(updatedOutputs[hex.EncodeToString(tx.ID)], &out)
			}
		}

		//GO through the new stuff and add then all in
		var err error
		for txID, outputs := range updatedOutputs {
			key, _ := hex.DecodeString(txID)
			err = b.Put(key, TransactionCTX.SerializeOutputs(outputs))
		}

		return err

	})

	if err != nil {
		log.Panic(err)
	}

}

func (UTXOSet *UTXO) FindTx(txId []byte, vout int) *TransactionCTX.TxOutput {

	var out *TransactionCTX.TxOutput
	err := UTXOSet.Blockchain.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(utxobucket))
		byteOutputs := b.Get(txId)

		outs := TransactionCTX.DeserializeTxOutput(byteOutputs)

		for _, txOut := range outs {
			if txOut.Index == vout {
				out = txOut
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
		return nil
	}
	return out
}

func deleteVout(txOutputs []*TransactionCTX.TxOutput, outputIdx int) []*TransactionCTX.TxOutput {

	for idx, vout := range txOutputs {
		if vout.Index == outputIdx {
			outputIdx = idx
			break
		}
	}
	txOutputs = append(txOutputs[:outputIdx], txOutputs[outputIdx+1:]...)
	return txOutputs
}

func NewUTXO(Blockchain *Blockchain) *UTXO {
	return &UTXO{Blockchain}
}
