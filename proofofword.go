package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 12 //前20位为0，用16进制表示，则前5个为0

type ProofOfWork struct {
	block  *Block
	target *big.Int //target为任意长度的int数
}

func newProofOfWord(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits)) //通过移位得到目标

	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)), //为啥要添加这个
			IntToHex(int64(nonce)),
		}, []byte{},
	)
	return data
}

/**
运行pow，知道计算出符合难度的hash
*/
func (pow *ProofOfWork) Run() (int, []byte) {
	nonce := 0
	var hash [32]byte
	var hashInt big.Int
	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		fmt.Printf("\r%d", nonce)

		hashInt.SetBytes(hash[:]) //将byte转为big.Int
		//Cmp()就是对整型指针进行操作的
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isVaild := (hashInt.Cmp(pow.target) == -1) //算出的hash值比目标小，isVaild为true
	return isVaild
}
