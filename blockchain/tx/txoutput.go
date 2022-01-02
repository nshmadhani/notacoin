package tx

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	"nshmadhani.com/notcoin/blockchain/util"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
	Index      int
}

func (out *TxOutput) Lock(address []byte) {
	out.PubKeyHash = util.PubHashFromAddress(address)
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func (out *TxOutput) String() string {
	str := ""
	str += fmt.Sprintf("Value = %d\n", out.Value)
	str += fmt.Sprintf("Index = %d\n", out.Index)

	str += fmt.Sprintf("PubKeyHash= %s\n", hex.EncodeToString(out.PubKeyHash))

	return str
}

func (out *TxOutput) Bytes() []byte {

	b := util.ToByteArray(int32(out.Value))
	b = append(b, out.PubKeyHash...)
	b = append(b, util.ToByteArray(int32(out.Index))...)
	return b
}

func SerializeOutputs(outs []*TxOutput) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(outs)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}
func DeserializeTxOutput(blockBytes []byte) []*TxOutput {

	var block []*TxOutput
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return block
}

func NewTxOuput(value int, address string, index int) TxOutput {
	out := TxOutput{value, nil, index}
	out.Lock([]byte(address))
	return out
}

//
