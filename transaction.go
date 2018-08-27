package main

import (
	"encoding/gob"
	"bytes"
	"log"
	"crypto/sha256"
	"fmt"
	"encoding/hex"
)
const subsidy = 10

type Transaction struct {
	ID []byte
	Vin []TXInput
	Vout []TXOutput
}

// IsCoinbase checks whether the transaction is coinbase
func(tx Transaction) IsCoinbase()  bool{
	return len(tx.Vin)==1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func(tx *Transaction) setID()  {
	//将一个结构体变量转为byte的操作
	var encoded bytes.Buffer
	var hash [32]byte
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx) //将tx变为字节存在Buffer结体中
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes()) //调用Buffer结构体中的Bytes()
	tx.ID = hash[:]
}

/**
创建挖矿交易，data可能传进来创世块的那句话
(挖矿得到的区块交易就是coinbase交易，该交易输入指向空，只有输出)
 */
func NewCoinbaseTX(to,data string)  *Transaction{

	if data == ""{
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{},-1,data} //ScriptSig 传data？
	txout := TXOutput{subsidy,to} //ScriptPubKey 传接受者地址？
	return &Transaction{nil,[]TXInput{txin},[]TXOutput{txout}}
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from,to string,amount int,bc *BlockChain) *Transaction{
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	//找零，生成能被支付方解锁(spend)的output
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.setID()

	return &tx
}

type TXInput struct {
	Txid []byte  //存储引用output所在的transaction的ID
	Vout int   	//存储引用的output在transaction的TXOutput数组中下标值
	ScriptSig string
}
//价值被存在TXOutput 中，TXInput只是对输出的引用
type TXOutput struct {
	Value int          //又因为将Value命名为value导致错误
	ScriptPubKey string
}

//检查input能否被提供的。。。与下一个函数到底什么区别？没理解
func(in *TXInput) CanUnlockOutputWith(unlockingData string) bool{
	return (in.ScriptSig == unlockingData)
}
//检查output能否被提供的秘钥（此时为简单字符串）解锁
func(out TXOutput) CanBeUnlockedWith(unlockingData string)  bool{
	return (out.ScriptPubKey == unlockingData)
}

