package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

/**
将整型转为字节型
*/
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}


// ReverseBytes reverses a byte array
//前后对应位置交换，可直接进行交换，不需要中间变量
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}