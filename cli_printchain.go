package main

import (
	"fmt"
	"strconv"
)

/**
通过迭代器打印区块链信息
*/
func (cli *CLI) printChain() {

	bc := NewBlockchain()
	defer bc.db.Close()
	bi := bc.Iterator()

	for {
		block := bi.next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.HashTransactions())
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := newProofOfWord(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

}
