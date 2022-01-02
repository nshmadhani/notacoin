package util

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"net"
	"os"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

/***************************************************************************************************
*									UTIL FUNCTIONS												****
*****************************************************************************************************/

func Base58Encode(addressHash []byte) string {
	return base58.Encode(addressHash)
}

func Base58Decode(address string) []byte {
	return base58.Decode(address)
}
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, _ := ecdsa.GenerateKey(curve, rand.Reader)
	pubkey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubkey
}
func HashPubKey(pubkey []byte) []byte {
	pub256 := sha256.Sum256(pubkey)

	ripemdHasher := ripemd160.New()
	_, _ = ripemdHasher.Write(pub256[:])
	return ripemdHasher.Sum(nil)
}

func PubHashFromAddress(address []byte) []byte {
	pubKeyHash := Base58Decode(string(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	return pubKeyHash
}

func FileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		// path/to/whatever exists
		return true
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		return false
	}
	return false
}
func Remove(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func GobDecode(encodedData []byte, data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func ToAAddr(address string) *net.TCPAddr {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Panic(err)
	}
	return addr
}

func ToByteArray(i int32) (arr []byte) {
	arr = make([]byte, 4)
	binary.BigEndian.PutUint32(arr[0:4], uint32(i))
	return
}
