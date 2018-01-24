package blockchain

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/asn1"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"log"
	"math"
	"time"
)

//Block of data
/*
	Block Structure
	Version: The version the block was created on
	Command: The command executed on the client machine
	Results: The results of an executed command
	Date: Time stamp of command execution
	PrivousHash: Hash of the previous block
	Nonce: Number used for mining
	Hash: Hash results of the block
		Command+Results+Date+PreviousHash+Signature+Nonce(8 Byte BE)
	Signature: Signature of a block from a private key
		Command+Results+Date+PreviousHash+Nonce(8 Byte BE)
*/
type Block struct {
	Version      float32 `json:"version"`
	Command      string  `json:"command"`
	Results      string  `json:"results"`
	TimeStamp    int64   `json:"time_stamp"`
	Target       string  `json:"target"`
	Nonce        uint64  `json:"nonce"`
	PreviousHash string  `json:"previous_hash"`
	Signature    string  `json:"signature"`
	Hash         string  `json:"hash"`
}

//Chain of blocks
type Chain struct {
	Blocks []Block `json:"blocks"`
	Keys
	*difficulty
}

//Keys for signing and validating blocks
type Keys struct {
	Private *rsa.PrivateKey
	Public  rsa.PublicKey
}

//Generate a RSA key pair for block signing
func (chain *Chain) setKeys() {
	keys, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	chain.Keys.Private = keys
	chain.Keys.Public = keys.PublicKey
}

//Generate random bytes for the Gensis block
func getRandomBytes(length int) []byte {
	tmp := make([]byte, length)
	_, err := rand.Read(tmp)
	if err != nil {
		log.Fatal(err)
	}
	return tmp
}

//ExportPEMKey helps mkae the public key more managable
func (chain *Chain) ExportPEMKey(key rsa.PublicKey) []byte {
	ans1Bytes, err := asn1.Marshal(key)
	if err != nil {
		log.Fatal(err)
	}
	pemKey := &pem.Block{Type: "PUBLIC KEY", Bytes: ans1Bytes}
	return pem.EncodeToMemory(pemKey)
}

//New block chain
func New() *Chain {
	block := &Block{}
	chain := &Chain{}
	chain.difficulty = initializeTarget()
	chain.difficulty.prvTime = time.Now().Unix()
	chain.setKeys()
	block.Version = 0.1
	block.TimeStamp = time.Now().Unix()
	block.Command = "Genesis Command " + hex.EncodeToString(getRandomBytes(64))
	block.Results = "Genesis Results " + hex.EncodeToString(getRandomBytes(64))
	pubBytes, err := asn1.Marshal(chain.Keys.Public)
	if err != nil {
		log.Fatal(err)
	}
	block.PreviousHash = hex.EncodeToString(pubBytes)
	chain.Mine(block)
	chain.Add(block)
	return chain
}

//Add block to the chain
func (chain *Chain) Add(block *Block) {
	chain.Blocks = append(chain.Blocks, *block)
	chain.difficulty.calcMineTime(block.TimeStamp)
}

//Create a sha512 hash of a pasted in byte array
func createHash(data []byte) []byte {
	hash := sha512.New()
	hash.Write(data)
	return hash.Sum(nil)
}

//Hash a byte array with a passed in nonce
func hashBlock(data []byte, nonce uint64) []byte {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	copyData := make([]byte, len(data))
	copy(copyData, data)
	copyData = append(copyData, nonceBytes...)
	return createHash(copyData)
}

//Take a block object and turn the fields into a byte array
func (chain *Chain) blockToBytes(block *Block) []byte {
	var data []byte
	var floatBytes [4]byte
	binary.BigEndian.PutUint32(floatBytes[:], math.Float32bits(block.Version))
	data = append(data, floatBytes[:]...)
	data = append(data, []byte(block.Command)...)
	data = append(data, []byte(block.Results)...)
	hashBytes, err := hex.DecodeString(block.PreviousHash)
	if err != nil {
		log.Fatal(err)
	}
	data = append(data, hashBytes...)
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, block.TimeStamp)
	data = append(data, buf...)
	data = append(data, []byte(block.Target)...)
	h := sha512.New()
	h.Write(data)
	hd := h.Sum(nil)
	signedData, err := chain.Keys.Private.Sign(rand.Reader, hd, crypto.SHA512)
	if err != nil {
		log.Fatal(err)
	}
	data = append(data, signedData...)
	block.Signature = hex.EncodeToString(signedData)
	return data
}
