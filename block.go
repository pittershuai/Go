package main

import (
	"time"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Data          []byte
	Hash          []byte
	nonce         int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), prevBlockHash, []byte(data),
		[]byte{}, 0}
	pow := newProofOfWord(block)
	nonce, hash := pow.Run() //删除setHash()通过pow计算出hash
	block.Hash = hash
	block.nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
