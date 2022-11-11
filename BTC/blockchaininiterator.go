package main

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	db *bolt.DB
	currentHashPointer []byte
}

func (bc *BlockChain) NewIterator() *BlockChainIterator  {
	return &BlockChainIterator{
		bc.db,
		bc.tail,
	}
}
//迭代器是属于区块链的，Next是属于迭代器的，返回当前区块，指针前移

func (it *BlockChainIterator ) Next()  *Block{
	var block Block
	_=it.db.View(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte(blockBucket))
//		fmt.Println(it)
//		fmt.Println(bucket)
		if bucket==nil{
			log.Panic("迭代器遍历时bucket不应该为空，请检查")
		}
//		fmt.Println(it.currentHashPointer)
		blockTmp :=bucket.Get(it.currentHashPointer)
//		fmt.Println(blockTmp)
		block=Deserialize(blockTmp)
		it.currentHashPointer=block.PrevHash
		return nil
	})
	return &block
}