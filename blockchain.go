package main

import (
	"fmt"
	"github.com/boltdb/bolt" //go get... 安装， 会下载到user/ss/go/bin目录下，添加这个路径到gopath即可
	"log"
)

const dbFile = "blockchain.db" //数据库文件
const blocksBucket = "blocks"  //数据库bucket的名称

//只需要保存数据库链接和最后一个block的hash即可
type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err = b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

//创建区块链。该方法是从无到有创建区块链。创建创世区块
func NewBlockchain() *BlockChain {
	//数据库中只存两类数据1、(block.hash-> block) 2、(l -> the hash of last block)

	var tip []byte                          //只向最后一个区块的Hash
	db, err := bolt.Open(dbFile, 0600, nil) //得到数据库链接

	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket)) //得到数据库中的名为blocks的bucket

		//若不存在该数据库，创建数据库-》创建创世块，并添加至数据库-》
		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")

			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize()) //创世块的hash作为key，创世块序列化作为value
			if err != nil {
				log.Panic(err)
			}
			//key l 对应val为区块链的最后一个区块的hash
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			//若存在数据库，直接返回blockchain最后一个block的hash
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := &BlockChain{tip, db}
	return bc
}
