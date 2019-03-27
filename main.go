package main

import "fmt"

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ismail")
	bc.AddBlock("Send 4 BTC to Ismail")

	for i, block := range bc.blocks {
		fmt.Printf("Prev Hash -- %x\n", block.PrevBlockHash)
		fmt.Printf("Data -- %s\n", block.Data)
		fmt.Printf("Hash -- %x", block.Hash)
		fmt.Printf("Index -- %d", i)
	}
}
