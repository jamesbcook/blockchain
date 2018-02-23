package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"math/big"
)

//ValidateHash of block
func (chain *Chain) ValidateHash(hash string, block *Block) (bool, error) {
	data := blockToBytes(block)
	sigBytes, err := hex.DecodeString(block.Signature)
	if err != nil {
		log.Fatal(err)
	}
	data = append(data, sigBytes...)
	res := hashBlock(data, block.Nonce)
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
	switch bytes.Compare(data, chain.difficulty.target) {
	case -1:
		return true
	case 1:
		return false
	}
	return false
}

//ValidateSig of block
func (chain *Chain) ValidateSig(sig string, block *Block) error {
	data := blockToBytes(block)
	h := sha256.New()
	h.Write(data)
	decodedSig, err := hex.DecodeString(sig)
	if err != nil {
		return err
	}
	r := big.Int{}
	s := big.Int{}
	sigLen := len(decodedSig)
	r.SetBytes(decodedSig[:(sigLen / 2)])
	s.SetBytes(decodedSig[(sigLen / 2):])
	x := big.Int{}
	y := big.Int{}
	keyLen := len(chain.Keys.Public)
	x.SetBytes(chain.Keys.Public[:(keyLen / 2)])
	y.SetBytes(chain.Keys.Public[(keyLen / 2):])
	rawPublic := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}
	if ecdsa.Verify(&rawPublic, h.Sum(nil), &r, &s) == false {
		return errors.New("Sig doesn't match")
	}
	return nil
}
