package main

import (
	"bytes"
)

type TXInput struct {
	Txid 		[]byte  //存储引用output所在的transaction的ID
	Vout 		int   	//存储引用的output在transaction的TXOutput数组中下标值
	Signature 	[]byte
	PubKey 		[]byte
}


/**
该input体能否引用pubKeyHash所在的output.等于解锁。上锁在TXOutput结构体的lock()中
 */
func(in *TXInput) UsesKey(pubKeyHash []byte)  bool{
	lockinghash := HashPubKey(in.PubKey) //
	return (bytes.Compare(lockinghash,pubKeyHash) == 0)
}


