package main

import "fmt"

func main() {
	bc := NewBlockchain()
	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 BTC to Ivan")

	//对于数组range返回（下标值，值）
	for _, block := range bc.blocks {
		pow := newProofOfWord(block)
		fmt.Println(pow.Validate())
	}
}
