package main

import (
	"fmt"
	"strconv"
)

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ismail")
	bc.AddBlock("Send 4 BTC to Ismail")

	for i, block := range bc.blocks {
		fmt.Printf("Prev Hash -- %x\n", block.PrevBlockHash)
		fmt.Printf("Data -- %s\n", block.Data)
		fmt.Printf("Hash -- %x\n", block.Hash)
		fmt.Printf("Index -- %d\n", i)

		pow := NewProofOfWork(block)
		fmt.Printf("POW: %s\n\n", strconv.FormatBool(pow.Validate()))
	}
}
