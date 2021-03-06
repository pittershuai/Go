package main

import (
	"fmt"
	"github.com/boltdb/bolt" //go get... 安装， 会下载到user/ss/go/bin目录下，添加这个路径到gopath即可
	"log"
	"os"
	"encoding/hex"
	"bytes"
	"errors"
	"crypto/ecdsa"
)

const dbFile = "blockchain.db" //数据库文件
const blocksBucket = "blocks"  //数据库bucket的名称
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"


//只需要保存数据库链接和最后一个block的hash即可
type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

/**
用于输出区块链中区块的迭代器结构体
当前区块的hash：用于在数据中找到结构体
数据库链接:需要在数据库中用hash找block，所以需要这个链接
*/
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

/**
往blockchain中添加block
由以前的addBlock改名为MineBlock
 */
func(bc *BlockChain) MineBlock(transactions []*Transaction) *Block{
	var lastHash []byte

	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
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

	return newBlock
}

// FindTransaction finds a transaction by its ID
func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.next()

		for _, tx := range block.Tracsactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

/**
returns a list of transactions containing unspent outputs
通过全部遍历的方式，遍历每个block中的每个Transaction中的每个TXOutput 和 TXInput。
只要发现Output被引用过（spend）就跳过循环。
*/
func(bc *BlockChain) FindUnspentTransactions(pubKeyHash []byte)  []Transaction{
	bci := bc.Iterator()
	var unspentTXs []Transaction
	//键为：被引用过的交易的hash（hash被转为string），
	//值为数组：该交易中被引用过的output的数组下标
	spentTXOs := make(map[string][]int)

	for{

		block := bci.next()

		for _,tx := range block.Tracsactions{
			txID := hex.EncodeToString(tx.ID) //将byte数组转为string类型
		Outputs:
			for outIdx,out := range tx.Vout{
				// 该交易中是否有Output被引用（spend）过。 这个判断能进去吗？spentTXOs[txID]没有被赋值应该一直为空啊。
				if spentTXOs[txID] != nil{
					//通过遍历被引用过的交易中被引用过的output，若发现spentTXOs[txID]有该记录则跳过
					for _,spentOut := range spentTXOs[txID]{
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			//只要不是coinbase就一定有input，记录到spentTXOs中
			if tx.IsCoinbase() == false{
				for _,in := range tx.Vin{
					if in.UsesKey(pubKeyHash){
						intxID := hex.EncodeToString(in.Txid)
						spentTXOs[intxID] = append(spentTXOs[intxID], in.Vout)
					}
				}
			}
		}


		if len(block.PrevBlockHash) == 0{ //PrevBlockHash为byte数组
			break
		}
	}
	return unspentTXs
}

// finds and returns unspent outputs to reference in inputs
//找到某地址上足够数目的UTXO即可停止，用于支付时调用吧
//返回支付的金额（包含找零的那部分），已经引用到哪些交易中的那些output
func (bc *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
//方法为寻找整个blockchain中的所有TUXO，UTXO_set.go中为寻找指定address下的UTXO。
func (bc *BlockChain) FindUTXO() map[string]TXOutputs {
	//key:transaction的id（byte[] 转为 string）
	//value:transation中未spend过的全部output
	UTXO := make(map[string]TXOutputs)
	//key:同上
	//value:transaton 中Vout数组的index
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.next()

		//按照此逻辑，最后一个transaction中的output一定是没有被引用过的。
		//因为每个block中只包含一个transaction。
		for _, tx := range block.Tracsactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			//只要是非coinbase的transaction，其他的tx都会有input
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					//将被引用过的transaction的output记录在spentTXOs中
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}



/**
放在BlockChain结构体中，创建迭代器。因为每个迭代器应该依附于blockchain。
即每个blockchain创建的时候都应该有含有一个迭代器
*/
func (bc *BlockChain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

// SignTransaction signs inputs of a Transaction
func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		fmt.Println("++++++++++++++++++++++++++++")
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction verifies transaction input signatures
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	//若为coinbase不用进行验证交易。
	if tx.IsCoinbase(){
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}


/**

 */
func (i *BlockchainIterator) next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.currentHash = block.PrevBlockHash

	return block
}


//为什么同时存在 CreateBlockchain()这个方法？
//CreateBlockchain()用于创世。
// NewBlockchain可以理解为在程序中得到已经存在于bolt数据库中的blockchain
func NewBlockchain() *BlockChain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

/**
CreateBlockchain creates a new blockchain DB
该方法只能调用一次，创建了一个区块链后不能再创建
 */
func CreateBlockchain(address string) *BlockChain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}