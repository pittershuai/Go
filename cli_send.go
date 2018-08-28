package main

import "fmt"

/*
添加了交易后，将以前的简单AddBlock()替换为MineBlock().
mine只有在进行交易时进行,此时只是每个区块中含有一个交易
 */
func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockchain()
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}
