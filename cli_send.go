package main

import (
	"fmt"
	"log"
)

/*
添加了交易后，将以前的简单AddBlock()替换为MineBlock().
mine只有在进行交易时进行,此时只是每个区块中含有一个交易
 */
func (cli *CLI) send(from, to string, amount int) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}
	bc := NewBlockchain()
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	//由发送者进行挖矿
	cbTx := NewCoinbaseTX(from, "")
	//一个send操作同时产生两个tx，一个为交易本身的tx，一个为挖矿生成的tx
	txs := []*Transaction{cbTx,tx}
	newBlock := bc.MineBlock(txs)
	UTXOSet.Update(newBlock)
	fmt.Println("Success!")
}
