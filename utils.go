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
