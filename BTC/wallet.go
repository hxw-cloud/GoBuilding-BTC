package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

//这里的钱包是一个结构，每一个钱包保存了公钥，私钥对

type Wallet struct {
	Private *ecdsa.PrivateKey
	//PubKey,这里的不存储原始的公钥，而是存储X和Y拼接的字符串，在校验的时候重新拆分（参考R,S传递）
	PubKey []byte
}

//创建钱包

func NewWallet() *Wallet {
	curve :=elliptic.P256()
	//生成私钥
	privateKey,err:=ecdsa.GenerateKey(curve,rand.Reader)
	if err!=nil{
		log.Panic("")
	}
	//生成公钥
	pubKeyOrig :=privateKey.PublicKey
	//拼接X,Y
	pubKey :=append(pubKeyOrig.X.Bytes(),pubKeyOrig.Y.Bytes()...)
	return &Wallet{Private: privateKey,PubKey: pubKey}
}


//生成地址

func (w *Wallet ) NewAddress() string {
	pubKey :=w.PubKey

	rip160HashValue :=HashPubKey(pubKey)

	version :=byte(00)
	payload :=append([]byte{version},rip160HashValue...)
	//checksum
	checkCode :=CheckSum(payload)
	payload =append(payload,checkCode...)

	address :=base58.Encode(payload)
	return address
}

func HashPubKey(data []byte) []byte {
	hash :=sha256.Sum256(data)
	//理解为编码器
	rip160hasher := ripemd160.New()
	_,err := rip160hasher.Write(hash[:])
	if err!=nil{
		log.Panic(err)
	}
	//返回rip160的hash结果
	rip160HashValue := rip160hasher.Sum(nil)
	return rip160HashValue
}

func CheckSum(data []byte) []byte {
	hash1:=sha256.Sum256(data)
	hash2 :=sha256.Sum256(hash1[:])
	//前四个字节校验码
	checkCode := hash2[:4]
	//25字节数据
	return checkCode
}

func IsValiAddress(address string) bool {
	//解码
	addressByte :=base58.Decode(address)
	if len(addressByte)<4 {
		return false
	}
	//取数据
	payload := addressByte[:len(addressByte)-4]
	checksum1 := addressByte[len(addressByte)-4: ]
	//做checksum函数
	checksum2 :=CheckSum(payload)

//	fmt.Printf("checksum1 : %x\n",checksum1)
//	fmt.Printf("checksum2 : %x\n",checksum2)
	//比较
	return bytes.Equal(checksum1,checksum2)
}