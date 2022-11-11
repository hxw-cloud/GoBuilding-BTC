package main

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const  blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"
//引入区块链

func NewBlockChain(address string) *BlockChain {
	//创建一个创世块，作为第一个区块添加到区块链中
//	genesisBlock:=GenesisBlock()
//	return&BlockChain{
//		blocks: []*Block{genesisBlock},
//		}
	var lastHash []byte
	db,err:=bolt.Open(blockChainDb,0600,nil)
	if err!=nil {
		log.Panic(err)
	}
	_=db.Update(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte(blockBucket))
		if bucket==nil{
			//创建
			bucket,err=tx.CreateBucket([]byte(blockBucket))
			if err!=nil{
				log.Panic("创建数据库失败")
			}
			genesisBlock:=GenesisBlock(address)
			//hash作为key,block作为字节流
			_=bucket.Put(genesisBlock.Hash,genesisBlock.Serialize())
			_=bucket.Put([]byte("LastHashKey"),genesisBlock.Hash)
			lastHash=genesisBlock.Hash
//			//测试
//			blockBytes:= bucket.Get(genesisBlock.Hash)
//			block:=Deserialize(blockBytes)
//			fmt.Printf("block info : %v\n",block)
		}else {
			//写入内容
			lastHash=bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})
	return &BlockChain{db,lastHash}
}
//创世块

func GenesisBlock(address string) *Block{
	coinbase := NewCoinbaseTX(address,"Go-创世块！")
	return NewBlock([]*Transacation{coinbase},[]byte{})
}

//重构

type BlockChain struct {
//	blocks []*Block
	db *bolt.DB
	tail []byte  //存储最后一个区块的hash
}
//添加区块

func (bc *BlockChain)AddBlock(txs []*Transacation)  {
	for _,tx := range txs {
		if !bc.VerifyTransacation(tx){
			fmt.Println("发现无效交易！")
			return
		}
	}
//	lastBlock:=bc.blocks[len(bc.blocks)-1]
//	prevHash :=lastBlock.Hash
//	//创建新的区块链
	db:=bc.db
	lastHash:=bc.tail
	_=db.Update(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte(blockBucket))
		if bucket==nil{
			log.Panic("请检查")
		}
		block :=NewBlock(txs,lastHash)
		_=bucket.Put(block.Hash,block.Serialize())
		_=bucket.Put([]byte("LastHashKey"),block.Hash)
		bc.tail=block.Hash
		return nil
	})
//	block:=NewBlock(data,prevHash)
//	//添加到区块链中
//	bc.blocks=append(bc.blocks,block)
}
//

//func (bc *BlockChain)PrintChain()  {
//	blockHeight :=0
//	bc.db.View(func(tx *bolt.Tx) error {
//		b:=tx.Bucket([]byte("blockBucket"))
//		b.ForEach(func(k, v []byte) error {
//			if bytes.Equal(k,[]byte("LastHashKey")){
//				return nil
//			}
//			block :=Deserialize(v)
//			fmt.Println("=========当前高度为",blockHeight,"=========")
//			blockHeight++
//			fmt.Printf("版本号： %d\n",block.Version)
//			fmt.Printf("前区块Hash值： %x\n",block.PrevHash)
//			fmt.Printf("梅克尔根： %x\n",block.MerkelRoot)
//			fmt.Printf("当前区块Hash值：%x\n",block.Hash)
//			fmt.Printf("时间搓： %d\n",block.TimeStamp)
//			fmt.Printf("难度值（未定义，随便写的）： %d\n",block.Difficulty)
//			fmt.Printf("随机数： %d\n",block.Nonce)
//			fmt.Printf("区块数据：%s\n",block.Transacations[0].TXInputs[0].PubKey)
//			return nil
//		})
//		return 	nil
//	})
//}

func (bc *BlockChain) FindUTXOs(PubKeyHash []byte) []TXOutput {
	var UTXO []TXOutput
	txs :=bc.FindUTXOTransactions(PubKeyHash)
	for _,tx:=range txs{
		for _,output :=range tx.TXOutputs{
			if bytes.Equal(PubKeyHash,output.PubKeyHash){
				UTXO=append(UTXO,output)
			}
		}
	}
	return UTXO
}


func (bc *BlockChain ) FindNeedUTXOs(senderPubKeyHash []byte ,amount float64) (map[string][]int64,float64) {
	var utxos =make(map[string][]int64)
	var calc float64

	txs :=bc.FindUTXOTransactions(senderPubKeyHash)
	for _,tx:=range txs{
		for i,output :=range tx.TXOutputs{
//			if from ==output.PubKeyHash{
			if bytes.Equal(senderPubKeyHash,output.PubKeyHash){
				//UTXO=append(UTXO,output)
				//实现逻辑的位置，找到自己需要的最少的utxo
				//把utxo加起来，统计一下当前utxo的总额，
				//比较是否满足转账需求，满足返回，不满足继续统计
				if calc <amount{
					utxos[string(tx.TXID)]=append(utxos[string(tx.TXID)],int64(i))
					calc+=output.Value
					if calc>=amount{
						return utxos,calc
						}
					}
				}
			}
		}
	return utxos,calc
}
func (bc *BlockChain) FindUTXOTransactions(senderPubKeyHash []byte) []*Transacation {
//	var UTXO []TXOutput
	var txs []*Transacation
	spentOutputs :=make(map[string][]int64)
	//遍历output，找到与之际相关的utxo
	//便利input，找到之际花费过的utxo
	it :=bc.NewIterator()
	for{
		//遍历区块
		block:=it.Next()
		//遍历交易
		for _,tx:=range block.Transacations{
			fmt.Printf("current txid : %x\n",tx.TXID)
			//遍历output，找到和自己相关的utxo（添加到output之前检查一下是否已经消耗掉）
			OUTPUT:
				for i,output :=range tx.TXOutputs {
					//在这里做一个过滤，将所有消耗过的outputs过滤和当前的所即将添加output对比一下
					//如果相同，这跳过，否者添加
					//如果当前交易id存在于我们已经标识的，则存在已消耗过的
					if spentOutputs[string(tx.TXID)] !=nil{
						for _,j:=range spentOutputs[string(tx.TXID)]{
							if int64(i)==j{
//								fmt.Println("这里跳过")
								continue OUTPUT
							}
						}
					}
//					if output.PubKeyHash==address
					if bytes.Equal(senderPubKeyHash,output.PubKeyHash){
//						UTXO=append(UTXO,output)
						//返回所有包含我相关的交易集合！！！！！！！！！！
						txs =append(txs,tx)
					}
				}
				//挖矿交易，直接跳过
				if !tx.IsCoinbase(){
					//遍历input，找到自己花费过的utxo集合
					for _,input :=range tx.TXInputs {
						//判断一下当前这个input和目标是否一致
//						if input.Sig ==address{
						pubKeyHash := HashPubKey(input.PubKey)
						if bytes.Equal(senderPubKeyHash,pubKeyHash){
							//indexArray :=spentOutputs[string(input.TXid)]
							//indexArray =append(indexArray,input.Index)
							spentOutputs[string(input.TXid)]=append(spentOutputs[string(input.TXid)],input.Index)
						}
					}
				}
		}
		if len(block.PrevHash)==0{
			fmt.Println("遍历完成退出")
			break
		}
	}
	return txs
}

func (bc *BlockChain ) FindTransactionByTXid(id []byte)  (Transacation,error){
	//遍历区块链
	//遍历交易
	//比较交易，找到直接退出
	//没找到返回空的transaction ,同时返回错误状态
	it :=bc.NewIterator()
	for  {
		block :=it.Next()
		//遍历交易
		for _,tx :=range block.Transacations{
			//比较交易，找到直接退出
			if bytes.Equal(  tx.TXID,id){
				return *tx,nil
			}
		}
		if len(block.PrevHash)==0{
			fmt.Printf("遍历区块链结束\n")
			break
		}
	}
	return Transacation{}, errors.New("无效的交易id，请检查！！！")
}

func (bc *BlockChain)  SignTransaction(tx *Transacation,privateKey *ecdsa.PrivateKey) {
	//签名，交易的最后签名
	prevTXs := make(map[string]Transacation)
	//找到所有的引用的交易
	//根据inputs来找，有多少input，就遍历多少此
	//找到目标交易，根据txid来找
	//添加到prevTXs
	for _,input :=range tx.TXInputs{
		//根据id查找交易本身，需要遍历整个区块链
		tx ,err:= bc.FindTransactionByTXid(input.TXid)
		if err!=nil{
			log.Panic("err : ",err)
		}
		prevTXs[string(input.TXid)] = tx
	}

	tx.Sign(privateKey,prevTXs)
}


func (bc *BlockChain ) VerifyTransacation(tx *Transacation) bool {
	//签名，交易的最后签名
	if tx.IsCoinbase(){
		return true
	}
	prevTXs := make(map[string]Transacation)
	//找到所有的引用的交易
	//根据inputs来找，有多少input，就遍历多少此
	//找到目标交易，根据txid来找
	//添加到prevTXs
	for _,input :=range tx.TXInputs{
		//根据id查找交易本身，需要遍历整个区块链
		tx ,err:= bc.FindTransactionByTXid(input.TXid)
		if err!=nil{
			log.Panic("err : ",err)
		}
		prevTXs[string(input.TXid)] = tx
	}

	return tx.Verify(prevTXs)

}












