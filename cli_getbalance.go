package main

import (
	"fmt"
	"log"
)

//获得账户余额
func (cli *CLI) getBalance(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := NewBlockchain()
	UTXOSet := UTXOSet{bc}
	//一直等到包含defer语句的函数执行完毕时，延迟函数（defer后的函数）才会被执行，
	// 而不管包含defer语句的函数是通过return的正常结束，还是由于panic导致的异常结束。
	defer bc.db.Close()

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)


	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
