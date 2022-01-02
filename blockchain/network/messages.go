package network

import (
	"bytes"
	"encoding/gob"
	"log"

	"nshmadhani.com/notcoin/blockchain/util"
)

type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func (v Version) Serialize() []byte {
	return append([]byte("Version000"), util.GobEncode(v)...)
}

func DeserializeVersion(blockBytes []byte) Version {
	var block Version
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))

	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}
	return block
}

type GetBlock struct {
	AddrFrom string
}

func (v GetBlock) Serialize() []byte {

	return append([]byte("GetBlock00"), util.GobEncode(v)...)
}

func DeserializeGetBlock(blockBytes []byte) GetBlock {
	var block GetBlock
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}
	return block
}

type Inv struct {
	Items    [][]byte
	Type     string
	AddrFrom string
}

func (v Inv) Serialize() []byte {

	return append([]byte("Inv0000000"), util.GobEncode(v)...)
}

func DeserializeInv(blockBytes []byte) Inv {
	var block Inv
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}
	return block
}

type GetData struct {
	Item     []byte
	Type     string
	AddrFrom string
}

func (v GetData) Serialize() []byte {

	return append([]byte("GetData000"), util.GobEncode(v)...)
}

func DeserializeGetData(blockBytes []byte) GetData {
	var block GetData
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}
	return block
}

type Block struct {
	Item     []byte
	AddrFrom string
}

func (v Block) Serialize() []byte {

	return append([]byte("Block00000"), util.GobEncode(v)...)
}

func DeserializeBlock(blockBytes []byte) Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}
	return block
}

type TX struct {
	Transaction []byte
	AddrFrom    string
}

func (v TX) Serialize() []byte {

	return append([]byte("TX00000000"), util.GobEncode(v)...)
}

func DeserializeTX(blockBytes []byte) TX {
	var block TX
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}
	return block
}

type Peers struct {
	Peers    []string
	AddrFrom string
}

func (v Peers) Serialize() []byte {
	return append([]byte("Peers00000"), util.GobEncode(v)...)
}

func DeserializePeers(blockBytes []byte) Peers {
	var block Peers
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}
	return block
}

type ReqPeers struct {
	AddrFrom string
}

func (v ReqPeers) Serialize() []byte {
	return append([]byte("ReqPeers00"), util.GobEncode(v)...)
}

func DeserializeReqPeers(blockBytes []byte) ReqPeers {
	var block ReqPeers
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}
	return block
}
