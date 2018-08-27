package main
//test git

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

//去除CLI结构体中的数据库链接db，解耦？数据库链接的事情全部交给blockchain
type CLI struct {}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

//func (cli *CLI) addBlock(data string) {
//	//cli.bc.AddBlock(data)  //替换为send中MineBlock
//}

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

func (cli *CLI) createBlockchain(address string) {
	bc := CreateBlockchain(address) //?不将此bc赋值到cli中？
	bc.db.Close()
	fmt.Println("Done!")

}

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

//获得账户余额
func (cli *CLI) getBalance(address string) {
	bc := NewBlockchain() //NewBlockchain()与CreateBlockchain()的区别
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}


func (cli *CLI) run() {
	cli.validateArgs()

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)



	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	//sendCmd 有两个参数，分别接受参数后数据
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
