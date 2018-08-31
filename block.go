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
	Tracsactions   []*Transaction
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
/**
对每个tx的id进行hash，（tx的id由tx本身经hash得到）
对应merkle树，不是对所有tx进行hash，而是对每个tx的hash进行拼接，再对此拼接值进行hash
这个函数为啥绑定在Block上？——看看POW中的调用就懂了
 */
func(b* Block) HashTransactions()  []byte{
	var transations [][]byte

	for _,tx := range b.Tracsactions{
		transations = append(transations,tx.ID)
	}
	mTree := NewMerkleTree(transations)

 	return mTree.RootNode.Data
}

func NewBlock(transations []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), prevBlockHash, transations,
		[]byte{}, 0}
	pow := newProofOfWord(block)
	nonce, hash := pow.Run() //删除setHash()通过pow计算出hash
	block.Hash = hash
	block.Nonce = nonce
	return block
}

/**
创建创世块交易，其中只包含一个交易
 */
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
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
