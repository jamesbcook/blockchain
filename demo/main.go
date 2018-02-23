package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jamesbcook/blockchain"
)

var (
	work    = make(chan blockchain.Block, 10)
	working = make(chan bool, 1)
)

func worker(chain *blockchain.Chain) {
	for {
		block := <-work
		block.PreviousHash = chain.Blocks[len(chain.Blocks)-1].Hash
		working <- true
		chain.Mine(&block)
		<-working
		chain.Add(&block)
	}
}

func printSpinner(chainDone chan bool) {
	spinner := []string{"\\", "-", "|", "-", "/", "-"}
	pos := 0
	for {
		select {
		case <-chainDone:
			fmt.Println("")
			return
		default:
			fmt.Printf("Building Genisis Block [%s]\r", spinner[pos])
			if pos < len(spinner)-1 {
				pos++
			} else {
				pos = 0
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func main() {
	chainDone := make(chan bool, 1)
	chain := &blockchain.Chain{}
	go func() {
		chain = blockchain.New()
		chainDone <- true
	}()
	printSpinner(chainDone)
	go worker(chain)
	prompt := "shell> "
	for {
		fmt.Printf("%s", prompt)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		if input == "" {
			continue
		} else if input == "exit" {
			waiting := 1
			for {
				fmt.Printf("Time Waiting... %d\r", waiting*2)
				if len(working) != 0 {
					waiting++
					time.Sleep(2 * time.Second)
				} else {
					println()
					break
				}
			}
			break
		}
		out, err := exec.Command("sh", "-c", input).Output()
		if err != nil {
			log.Println(err)
			out = []byte(err.Error())
		}
		stringOut := string(out)
		fmt.Println(stringOut)
		createBlock(input, stringOut[:len(stringOut)-1])
	}

	//Write the current chain to a log file
	err := chain.Write("blockchain.log")
	if err != nil {
		log.Fatal(err)
	}

	//Check an individual block has a valid hash
	valid, err := chain.ValidateHash(chain.Blocks[2].Hash, &chain.Blocks[2])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("hash valid?:", valid)

	//Check an individual block has a valid signature
	err = chain.ValidateSig(chain.Blocks[2].Signature, &chain.Blocks[2])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("sig valid?", true)

	//Loop over each block in the chain and confirm they have a valid
	//hash and signature.
	for x, block := range chain.Blocks {
		if _, err := chain.ValidateHash(block.Hash, &block); err != nil {
			log.Fatal(err)
		}
		if err := chain.ValidateSig(block.Signature, &block); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Block %d's hash and signature match\n", x)
	}
}

func createBlock(command, res string) {
	block := &blockchain.Block{}
	block.Command = command
	block.Results = res
	work <- *block
}
