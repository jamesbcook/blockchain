package blockchain

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"math"
)

//ValidateHash of block
func (chain *Chain) ValidateHash(hash string, block *Block) (bool, error) {
	data := chain.blockToBytes(block)
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, block.Nonce)
	data = append(data, nonceBytes...)
	res := createHash(data)
	decodeHash, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	if bytes.Compare(res, decodeHash) == 0 {
		return true, nil
	}
	return false, nil
}

//Validate the hash meets the required difficulty
func (chain *Chain) validHash(data []byte) bool {
	match := 0
	for i := 0; i < chain.Difficulty; i++ {
		if data[i] != 0 {
			break
		} else {
			match++
		}
		if match == chain.Difficulty {
			return true
		}
	}
	return false
}

//ValidateSig of block
func (chain *Chain) ValidateSig(sig string, block *Block) error {
	var data []byte
	var floatBytes [4]byte
	binary.BigEndian.PutUint32(floatBytes[:], math.Float32bits(block.Version))
	data = append(data, floatBytes[:]...)
	data = append(data, []byte(block.Command)...)
	data = append(data, []byte(block.Results)...)
	hashBytes, err := hex.DecodeString(block.PreviousHash)
	if err != nil {
		return err
	}
	data = append(data, hashBytes...)
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, block.TimeStamp)
	data = append(data, buf...)
	h := sha512.New()
	h.Write(data)
	decodedSig, err := hex.DecodeString(sig)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(&chain.Keys.Public, crypto.SHA512, h.Sum(nil), decodedSig)
}
