package blockchain

import (
	"encoding/hex"
	"sync"
)

const (
	threads = 12
)

//Mine a block
func (chain *Chain) Mine(block *Block) {
	type foundData struct {
		nonce uint64
		hash  string
	}
	dataChannel := make(chan foundData, threads)
	nonceChannel := make(chan uint64, threads*2)
	doneChannel := make(chan bool, threads)
	var wg sync.WaitGroup
	var nonce uint64
	found := false
	data := chain.blockToBytes(block)
	for x := 0; x < threads*2; x++ {
		nonceChannel <- nonce
		nonce++
	}
	for x := 0; x < threads; x++ {
		wg.Add(1)
		go func() {
			for {
				tmpNonce := <-nonceChannel
				res := hashBlock(data, tmpNonce)
				if chain.validHash(res) && found == false {
					found = true
					dataChannel <- foundData{nonce: tmpNonce,
						hash: hex.EncodeToString(res)}
					for i := 0; i < threads; i++ {
						doneChannel <- true
					}
				}
				select {
				case f, ok := <-doneChannel:
					if ok && f {
						wg.Done()
						return
					}
				default:
					nonce++
					nonceChannel <- nonce
				}
			}
		}()
	}
	wg.Wait()
	hashedBlock := <-dataChannel
	block.Nonce = hashedBlock.nonce
	block.Hash = hashedBlock.hash
	close(dataChannel)
	close(doneChannel)
	close(nonceChannel)
}
