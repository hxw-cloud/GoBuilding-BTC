package main

import (
	"bytes"
	"crypto/sha256"
	"math/big"
)

type ProofOfWork struct {
	//block
	block *Block
	//一个非常大的数
	target *big.Int
}

//提过pow函数


func NewProofOfWork(block *Block) *ProofOfWork{
	pow:=ProofOfWork{
		block: block,
	}
	targetStr:="0001000000000000000000000000000000000000000000000000000000000000"
	tmpInt:=big.Int{}
	tmpInt.SetString(targetStr,16)
	pow.target = &tmpInt
	return &pow
}
//提供不断计算hash的函数

func (pow *ProofOfWork)Run()([]byte,uint64)  {
	//拼装数据（区块数据，不断变化的随机数）
	//做hash运算
	//与pow中的target进行比较
	//找到了返回，没找到随机数加一
	var nonce uint64
	block:=pow.block
	var hash [32]byte
	for{
		tmp:=[][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(nonce),
//			block.Data,
//只对区块头做Hash,区块体通过Merkel root产生影响
			}
			//将二维的数组切片连接，返回一个一维切片
			blockInfo:=bytes.Join(tmp,[]byte{})
			hash=sha256.Sum256(blockInfo)
			tmpInt:=big.Int{}
			tmpInt.SetBytes(hash[:])
			if tmpInt.Cmp(pow.target)==-1{
				break
			}else{
				nonce++
			}
//			if nonce>10000{
//				fmt.Print("not find\n")
//				break
//			}
	}
	return hash[:],nonce
}
//