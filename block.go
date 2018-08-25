package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Data          []byte
	Hash          []byte
	Nonce         int //Nonce写为nonce导致错误（还是不报错的错误）,命名也是很重要的
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), prevBlockHash, []byte(data),
		[]byte{}, 0}
	pow := newProofOfWord(block)
	nonce, hash := pow.Run() //删除setHash()通过pow计算出hash
	block.Hash = hash
	block.Nonce = nonce
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
