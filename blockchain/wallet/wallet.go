package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"io/ioutil"
	"log"

	. "nshmadhani.com/notcoin/blockchain/tx"
	util "nshmadhani.com/notcoin/blockchain/util"
)

const version = "1"
const addressCheckSumLen = 4

type Wallet struct {
	PrivateKey        ecdsa.PrivateKey
	Publickey         []byte
	Alias             string
	EncodedPrivateKey []byte
}

func (w Wallet) GetAddress() string {
	pubKeyHash := util.HashPubKey(w.Publickey)
	versionPayload := append([]byte(version), pubKeyHash...)
	payload := append(versionPayload, w.Checksum()...)
	return util.Base58Encode(payload)
}

func (w Wallet) Checksum() []byte {
	f := sha256.Sum256(w.Publickey)
	s := sha256.Sum256(f[:])
	return s[:addressCheckSumLen]
}

func (w *Wallet) Serialize() []byte {
	var result bytes.Buffer
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(w)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeserializeWallet(walletBytes []byte) Wallet {

	var wallet Wallet
	gob.Register(elliptic.P256())

	decoder := gob.NewDecoder(bytes.NewReader(walletBytes))
	err := decoder.Decode(&wallet)

	if err != nil {
		log.Panic(err)
	}

	return wallet
}

func (w *Wallet) Sign(tx *Transaction, utxos map[string][]*TxOutput) *Transaction {
	if tx.IsCoinbase() {
		return tx
	}

	trimmedCopy := tx.Trim()

	for idx, in := range trimmedCopy.Vin {

		trimmedCopy.Vin[idx].Signature = nil
		utxo := utxos[hex.EncodeToString(in.TxId)]

		for _, out := range utxo {
			if out.Index == in.Vout {
				trimmedCopy.Vin[idx].PubKey = out.PubKeyHash
			}
		}
		trimmedCopy.SetId()
		trimmedCopy.Vin[idx].PubKey = nil

		r, s, _ := ecdsa.Sign(rand.Reader, &w.PrivateKey, trimmedCopy.ID)
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[idx].Signature = signature
	}

	return tx
}

/******************************************************************************************************
										Wallet Functions
*******************************************************************************************************/
func NewWallet(alias string) *Wallet {
	p, pub := util.NewKeyPair()
	return &Wallet{p, pub, alias, nil}
}

func GetWallet(walletFile string) Wallet {
	if !util.FileExists(walletFile) {
		aWallet := *NewWallet("RANDOM")
		_ = ioutil.WriteFile(walletFile, aWallet.Serialize(), 0644)
		return aWallet
	} else {
		data, _ := ioutil.ReadFile(walletFile)
		return DeserializeWallet(data)
	}
}
