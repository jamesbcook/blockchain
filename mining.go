package blockchain

import (
	"context"
	"encoding/hex"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var (
	threads = runtime.NumCPU()
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

//Mine a block
func (chain *Chain) Mine(block *Block) {
	type foundData struct {
		nonce uint64
		hash  string
	}
	dataChannel := make(chan foundData, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}
	block.Version = 0.1
	block.TimeStamp = time.Now().Unix()
	block.Target = hex.EncodeToString(chain.difficulty.targetBits)
	data := blockToBytes(block)
	sig := chain.signData(data)
	data = append(data, sig...)
	block.Signature = hex.EncodeToString(sig)
	for x := 0; x < threads; x++ {
		wg.Add(1)
		nonce := rand.Int63()
		go func(ctx context.Context, wg *sync.WaitGroup, nonce uint64) {
			defer wg.Done()
			for {
				res := hashBlock(data, nonce)
				if chain.validHash(res) {
					dataChannel <- foundData{nonce: nonce,
						hash: hex.EncodeToString(res)}
				}
				select {
				case <-ctx.Done():
					return
				default:
					nonce++
					if nonce >= math.MaxUint64 {
						nonce = 0
					}
				}
			}
		}(ctx, &wg, uint64(nonce))
	}
	hashedBlock := <-dataChannel
	block.Nonce = hashedBlock.nonce
	block.Hash = hashedBlock.hash
	cancel()
	wg.Wait()
}
