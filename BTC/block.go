package main
import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)
type Block struct{
	//1.版本号
	Version uint64
	//2.前驱hash
	PrevHash []byte
	//3.merkelroot 梅克尔根，hash值
	MerkelRoot []byte
	//4.时间搓
	TimeStamp uint64
	//5.难度值
	Difficulty uint64
	//6，random
	Nonce uint64
	//当前区块hash，正常的没有，方便简化
	Hash  []byte
	//数据
//	Data []byte
	Transacations []*Transacation
}
//创建区块

func NewBlock(txs []*Transacation,prevBlockHash []byte) *Block  {
	block:=Block{
		Version: 00,
		PrevHash: prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp: uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce: 0,
		Hash: []byte{},
//		Data: []byte(data),
		Transacations: txs,
		}
		block.MerkelRoot=block.MakeMerkelRoot()

//	block.SetHash()
	pow:=NewProofOfWork(&block)
	block.Hash,block.Nonce=pow.Run()
	return &block
}
//生成Hash
//func (block *Block)SetHash()  {
//	//1.拼装数据
//	var blockInfo []byte
////	blockInfo=append(blockInfo,Uint64ToByte(block.Version)...)
////	blockInfo=append(blockInfo,block.PrevHash...)
////	blockInfo=append(blockInfo,block.MerkelRoot...)
////	blockInfo=append(blockInfo,Uint64ToByte(block.TimeStamp)...)
////	blockInfo=append(blockInfo,Uint64ToByte(block.Difficulty)...)
////	blockInfo=append(blockInfo,Uint64ToByte(block.Nonce)...)
////	blockInfo=append(blockInfo,block.Data...)
//	tmp:=[][]byte{
//		Uint64ToByte(block.Version),
//		block.PrevHash,
//		block.MerkelRoot,
//		Uint64ToByte(block.TimeStamp),
//		Uint64ToByte(block.Difficulty),
//		Uint64ToByte(block.Nonce),
//		block.Data,
//	}
//	//将二维的数组切片连接，返回一个一维切片
//	blockInfo=bytes.Join(tmp,[]byte{})
//	//2.sha256
//	hash:=sha256.Sum256(blockInfo)
//	block.Hash=hash[:]
//}


//辅助函数，将uint转成byte

func Uint64ToByte(num uint64) []byte  {
	var buffer bytes.Buffer
	err:=binary.Write(&buffer,binary.BigEndian,num)
	if err!=nil{
		log.Panic(err)
	}
	return buffer.Bytes()
}
//序列化

func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder:=gob.NewEncoder(&buffer)
	err:=encoder.Encode(&block)
	if err!=nil{
		log.Panic("编码出错")
	}
	return buffer.Bytes()
}

//反序列化

func Deserialize(data []byte) Block {
    decoder:=gob.NewDecoder(bytes.NewReader(data))
	var block Block
	err:=decoder.Decode(&block)
	if err!=nil{
		log.Panic("解码出错",err)
	}
	return block
}


func (block *Block ) MakeMerkelRoot() []byte {
	//将交易的hash拼接起来，再整体做hash处理
	var info []byte
	for _,tx:=range block.Transacations{
		info=append(info,tx.TXID...)
	}
	hash :=sha256.Sum256(info)
	return hash[:]
}
