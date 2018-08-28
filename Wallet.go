package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha1"

	"golang.org/x/crypto/ripemd160"

	"crypto/sha256"
	"bytes"
)

const version = byte(0x00)  //版本
const addressChecksumLen = 4  //检查长度4字节
const walletFile = "wallet.dat"

/*
存放秘钥对
 */
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}
/**
生成过程参照由pubKey生成address的图,将过程分为到两个函数中
 */
func(w *Wallet) GetAddress()  []byte{
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha1.Sum(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

/**
ValidateAddress check if address if valid
将地址转换为三部分（version、pubKeyHash、checksum），再组合进行checksum，得到新的checksum，两个checksum进行比对
 */
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

/*
每个结构体都应该有相应的New方法
 */
func NewWallet()  *Wallet{
	private,public := NewKeyPair()
	return &Wallet{private,public}
}

func NewKeyPair()(ecdsa.PrivateKey,[]byte)  {
	curve := elliptic.P256()
	private,err := ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//...表示接受任意数量个参数
	pubKey := append(private.PublicKey.X.Bytes(),private.PublicKey.Y.Bytes()...)
	return  *private,pubKey
}


