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
	//挖矿得到的block中的TXInput是没有信息的，所以在PubKey中放点信息也是OK的。
	txin := TXInput{[]byte{},-1,nil,[]byte(data)}
	txout := NewTXOutput(subsidy,to)
	tx := Transaction{nil,[]TXInput{txin},[]TXOutput{*txout}}
	tx.setID()
	return &tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from,to string,amount int,bc *BlockChain) *Transaction{
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)

	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

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
			input := TXInput{txID, out, nil,wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, *NewTXOutput(amount,to))
	//找零，生成能被支付方解锁(spend)的output
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc - amount,from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.setID()

	return &tx
}