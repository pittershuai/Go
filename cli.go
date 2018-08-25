package main
//test git

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	bc *BlockChain
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}
func (cli *CLI) printUsage() {
	fmt.Println("Usage")
	fmt.Println(" addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println(" printchain - print all the blocks of the blockchain")
}
func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
}

/**
通过迭代器打印区块链信息
*/
func (cli *CLI) printChain() {
	bi := cli.bc.Iterator()
	for {
		block := bi.next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
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

func (cli *CLI) run() {
	cli.validateArgs()

	// NewFlagSet的第一个参数是可以任意定的.但第二个参数,则决定了参数解析出错时错误处理方式.
	//分为两步1：让命令行识别命令名称，2：得到命令名称后相应参数的值

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError) //可以理解为在命令行中注册了addblock，让命令行能识别
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	//将输入中blockchain_go addblock -data “。。。。” name为data，返回data后字符串的指针
	addBlockData := addBlockCmd.String("data", "", "Block data")

	//先运行 go build，生成blockchain_go.exe文件
	//再执行blockchain_go printchain 或 再执行blockchain_go addblock -data “。。。。”
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
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
	//是否使用addBlockCmd解析过参数
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData) //addBlockData,*为取值操作
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}

}
