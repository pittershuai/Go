package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

//价值被存在TXOutput 中，TXInput只是对输出的引用
type TXOutput struct {
	Value int          //又因为将Value命名为value导致错误
	PubKeyHash []byte
}

/*
给output上锁 (使用address所含有的pubKey，将该pubKey经hash后上锁)
 */
func(out *TXOutput) lock(address []byte)  {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1:len(pubKeyHash)-4] //版本号占1字节，checksum占4字节
	out.PubKeyHash = pubKeyHash
}

/**
确定是否output是否在真的被address的pubKey锁定。用于output被引用的时候确认是能被引用
 */
func(out *TXOutput) IsLockedWithKey(pubKeyHash []byte)  bool{
	return bytes.Compare(pubKeyHash,out.PubKeyHash) == 0
}

/*
传入字符串、价值，生成output
 */
func NewTXOutput(value int ,address string)  *TXOutput{
	txo := &TXOutput{value, nil}
	txo.lock([]byte(address))
	return txo
}

/**
为了装多个TXOutput，并进行序列化，反序列化。这个结构体有必要吗？
 */
type TXOutputs struct {
	Outputs []TXOutput
}

// Serialize serializes TXOutputs
func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs deserializes TXOutputs
func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}


