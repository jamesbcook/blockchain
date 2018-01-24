package blockchain

import (
	"encoding/binary"
	"fmt"
)

//Difficulty data container
type difficulty struct {
	targetBits  []byte
	target      []byte
	minedBlocks int
	prvTime     int64
	blockTime   int64
}

const (
	difficultyBlock = 128
)

//Number of seconds it should take to mine a block
func getTargetTime() float64 {
	return 10.0
}

func byteToInt(input []byte) uint64 {
	return binary.BigEndian.Uint64(input)
}

func intToByte(input uint64) []byte {
	tmpBuffer := make([]byte, 64)
	binary.BigEndian.PutUint64(tmpBuffer, input)
	return tmpBuffer
}

func (d *difficulty) calcMineTime(mineTime int64) {
	d.minedBlocks++
	d.blockTime += mineTime - d.prvTime
	d.prvTime = mineTime
	if d.minedBlocks == difficultyBlock {
		avgTime := float64(d.blockTime) / float64(difficultyBlock)
		fmt.Println(avgTime)
		targetTime := getTargetTime()
		tmpBuffer := make([]byte, 64)
		copy(tmpBuffer[d.targetBits[0]:], d.targetBits[1:])
		num := byteToInt(tmpBuffer)
		exponent := 0
		if avgTime < targetTime {
			newTargetnumber := float64(num) / (targetTime / avgTime)
			newTargetBytes := intToByte(uint64(newTargetnumber))
			for _, v := range newTargetBytes {
				if v == 0 {
					exponent++
				} else {
					break
				}
			}
			targetBitBytes := [8]byte{}
			targetBitLen := 1
			for x := exponent; x < 64; x++ {
				if newTargetBytes[x] != 0 {
					targetBitLen++
				} else {
					break
				}
			}
			targetBitBytes[0] = byte(exponent)
			copy(targetBitBytes[1:], newTargetBytes[exponent:])
			d.targetBits = targetBitBytes[:targetBitLen]
			d.target = newTargetBytes[:]
		} else if avgTime > targetTime {
			percentOver := 1.0 - (targetTime / avgTime)
			numToAdd := float64(num) * percentOver
			newTargetnumber := float64(num) + numToAdd
			newTargetBytes := intToByte(uint64(newTargetnumber))
			for _, v := range newTargetBytes {
				if v == 0 {
					exponent++
				} else {
					break
				}
			}
			targetBitBytes := [8]byte{}
			targetBitLen := 1
			for x := exponent; x < 64; x++ {
				if newTargetBytes[x] != 0 {
					targetBitLen++
				} else {
					break
				}
			}
			targetBitBytes[0] = byte(exponent)
			copy(targetBitBytes[1:], newTargetBytes[exponent:])
			d.targetBits = targetBitBytes[:targetBitLen]
			d.target = newTargetBytes[:]
		}
		fmt.Println(d.target)
		fmt.Println(d.targetBits)
		d.minedBlocks = 0
		d.blockTime = 0
	}
}

//InitializeTarget stuff
func initializeTarget() *difficulty {
	bits := []byte{0x02, 0xff, 0xff}
	dif := &difficulty{}
	dif.minedBlocks = 0
	dif.blockTime = 0
	temp := make([]byte, 64)
	copy(temp[bits[0]:], bits[1:])
	dif.targetBits = bits[:]
	dif.target = temp
	return dif
}
