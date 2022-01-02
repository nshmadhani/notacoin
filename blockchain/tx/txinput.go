package tx

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"nshmadhani.com/notcoin/blockchain/util"
)

type TxInput struct {
	Vout      int    // Index value of a tranactions txOut
	TxId      []byte // transcation from where we took it
	Signature []byte
	PubKey    []byte
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := util.HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (in *TxInput) String() string {
	str := ""
	str += fmt.Sprintf("Vout = %d\n", in.Vout)
	str += fmt.Sprintf("TxId = %s\n", hex.EncodeToString(in.TxId))
	str += fmt.Sprintf("Signature = %s\n", hex.EncodeToString(in.Signature))
	str += fmt.Sprintf("Pubkey = %s\n", hex.EncodeToString(in.PubKey))
	return str
}

func (in *TxInput) Bytes() []byte {

	b := util.ToByteArray(int32(in.Vout))
	b = append(b, in.TxId...)
	b = append(b, in.Signature...)
	b = append(b, in.PubKey...)

	return b
}
