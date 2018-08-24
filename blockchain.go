package main

type BlockChain struct {
	blocks []*Block //指針數組，否則直接用數組存放區塊會很占內存
}

//為BlockChain添加方法
func (bc *BlockChain) AddBlock(data string) {
	preBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, preBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

//创建区块链。该方法是从无到有创建区块链。创建创世区块
func NewBlockchain() *BlockChain {
	//为结构体赋值BlockChain{}，
	// 为结构体中数组赋值[]*Block{NewGenesisBlock(),...}
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}
