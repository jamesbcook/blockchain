package blockchain

import (
	"encoding/json"
	"io/ioutil"
)

//Load blockchain from a file
func (chain *Chain) Load(file string) (*Chain, error) {
	oldChain := &Chain{}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, oldChain)
	if err != nil {
		return nil, err
	}
	return oldChain, nil
}

//Write blockchain to a file
func (chain *Chain) Write(file string) error {
	content, err := json.MarshalIndent(chain.Blocks, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, content, 0644)
}
