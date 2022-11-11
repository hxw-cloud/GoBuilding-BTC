package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"strings"
)

//定义交易结构

type Transacation struct {
	TXID []byte
	TXInputs []TXInput
	TXOutputs []TXOutput
}
//定义交易输入
const reward = 50
type TXInput struct {
	TXid []byte
	Index int64
	//Sig string
	//真正的数子签名，由r,s拼接成城的【】byte
	Signature []byte
	PubKey []byte
}
//定义交易输出

type TXOutput struct {
	Value float64
/*	 string
	收款方的公钥的hash，注意是hash不是公钥，也不是地址*/
	PubKeyHash []byte
}
//由于现在存储的是公钥的hash，所以无法直接创建txoutput
//为了能够得到公钥hash，我们需要处理一下，写一个Lock函数

func (output *TXOutput ) Lock(address string)  {
//解码
	//截取出公钥hash：去除version（一字节），去除校验码（4字节）
//	addressByte := base58.Decode(address)  //25字节
//	lens := len(addressByte)
//	pubKeyHash := addressByte[1:lens-4]

	//真正的锁定动作
	output.PubKeyHash= GetPubKeyFromAddress(address)
}
//实现给txoutput提供一个创建的方法，否者无法调用lock

func NewTXOutput(value float64,address string)  *TXOutput {
	output := TXOutput{
		Value : value,
	}
	output.Lock(address)
	return &output
}


//设置交易id

func (tx *Transacation ) SetHash()  {
	var buffer bytes.Buffer
	encoder :=gob.NewEncoder(&buffer)
	err:=encoder.Encode(tx)
	if err!=nil{
		log.Panic(err)
	}
	data:=buffer.Bytes()
	hash :=sha256.Sum256(data)
	tx.TXID=hash[:]
}
//提供交易方法(挖矿交易)

func NewCoinbaseTX(address string,data string) *Transacation {
	//特点，只有一个input，无需引用交易id，无需引用index
	//签名先为空，最后完整交易后，最后做一次签名即可
	input:=TXInput{[]byte{},-1,nil,[]byte(data)}
	//新的创建方法
	output :=NewTXOutput(reward,address)
//	fmt.Println("output value:",output)
	transacation :=Transacation{[]byte{},[]TXInput{input},[]TXOutput{*output}}
//	fmt.Println("transacation : ",transacation)
	transacation.SetHash()
//	fmt.Println("transacation.sethash : ",transacation.TXID)
	return &transacation
}


//创建普通转账交易
//1.找到最合理的utxo集合，map[string][]int64
//创建outputs
//如果有零钱，要找零

func NewTransaction(from,to string,amount float64,bc *BlockChain) *Transacation{
	//创建交易之后要进行数子签名，->所以要私钥->打开钱包“NewWallets()”
	ws := NewWallets()
	//找到自己的钱包，根据地址返回自己的wallet
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Printf("没有找到改地址的钱包，交易创建失败!\n")
		return nil
	}
	//3 得到对应的公钥，私钥
	pubKey := wallet.PubKey
	privateKey := wallet.Private
	pubKeyHash := HashPubKey(pubKey)
	utxos,resValue := bc.FindNeedUTXOs(pubKeyHash,amount)
//	fmt.Println("resValue的值为:",resValue)
	if resValue<amount{
		fmt.Println("余额不足,所剩余额为：",resValue)
		return nil
	}
	var inputs []TXInput
	var outputs []TXOutput
	for id,indexArray:=range utxos{
		for _,i :=range indexArray {
			input :=TXInput{[]byte(id),i,nil,pubKey}
			inputs=append(inputs,input)
		}
	}
	output :=NewTXOutput(amount,to)
	outputs=append(outputs,*output)
	if resValue >=amount{
		output = NewTXOutput(resValue - amount,from)
		outputs=append(outputs,*output)
//		fmt.Println("TXOutput{resValue-amount,from} : ",TXOutput{resValue-amount,from})
//		fmt.Println("outputs : ",outputs)
	}
	tx:=Transacation{[]byte{},inputs,outputs}
	tx.SetHash()
	//测试使用

	bc.SignTransaction(&tx,privateKey)
//	fmt.Println("tx.TXOutputs[0].Value : ",tx.TXOutputs[0].Value)
	return &tx
}

//签名的具体实现,参数为：私钥，inputs里所以要引用的交易的结构map【string】Transacation
//map[2222]Transacation

func (tx *Transacation ) Sign(privateKey *ecdsa.PrivateKey ,prevTXs map[string]Transacation)  {

	if tx.IsCoinbase(){
		return
	}
	//创建一个当前交易copy ：txcopy，使用函数 TrimmedCopy ：要把Signature和PubKey字段设置为nil
	txCopy := tx.TrimmedCopy()
	//循环遍历txCopy的inputs，得到这个input的索引的output的公钥hash
	for i,input := range txCopy.TXInputs{
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID)==0{
			log.Panic("引用的交易无效")
		}
		//不要对input进行赋值，这是一个副本，要对txCopy.TXInput[xx]进行操作，否者无法把pubKeyHash传进来
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		/*
		所需要的三个数据都具备了，开始做hash处理
		生成要签名的数据，要签名的数据一定是hash值
		我们对每一个input都要签名一次，签名的数据是当前input引用的output的hash+当前的outputs(都承载在当前这个txCopy里面)
		要对这个拼好的txCopy进行hash处理，setHash得到的TXID,这个TXID就是我们要签名的最终数据
		*/
		txCopy.SetHash()
		//还原，以免影响后面的签名
		txCopy.TXInputs[i].PubKey =nil
		signDataHash := txCopy.TXID
		r,s,err := ecdsa.Sign(rand.Reader,privateKey,signDataHash)
		if err!=nil{
			log.Panic("err : ",err)
		}
		signature := append(r.Bytes(),s.Bytes()...)
		tx.TXInputs[i].Signature = signature
	}
	//生成要签名的数据。要签名的数据一定是hash值
		//我们要对每一个input都要签名一次，签名的数据是由当前input引用的output的hash+当前的outputs  （都承载在当前这个txCopy）
		//对这个品好的txCopy进行hash处理，setHash得到的TXID,这个TXID就是我们要签名的最终数据。
	//执行签名动作得到r,s字节流
	//放到我们所签名的input的Signature中

}

func (tx *Transacation ) TrimmedCopy() Transacation {
	var inputs []TXInput
	var outputs []TXOutput
	for _ ,input :=range tx.TXInputs{
		inputs = append(inputs,TXInput{input.TXid,input.Index,nil,nil})
	}
	for _,output := range tx.TXOutputs{
		outputs = append(outputs,output)
	}
	return  Transacation{tx.TXID,inputs,outputs}
}


//分析校验
//所需要的数据 : 公钥，数据（txCopy，生成hash），签名
//我们要对每一个签名过的input进行校验

func (tx *Transacation ) Verify(prevTXs map[string]Transacation) bool {
	if tx.IsCoinbase(){
		return true
	}
	//得到签名的数据
	txCopy := tx.TrimmedCopy()

	for i,input := range tx.TXInputs{
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID)  ==0{
			log.Panic("引用的交易无效")
		}
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		txCopy.SetHash()
		dataHash := txCopy.TXID
		//得到signature，反推r,s
		signature := input.Signature
		//拆解PubKey,X,Y得到原生公钥
		pubKey := input.PubKey //拆X,Y
		//定义两个辅助的big,int
		r:=big.Int{}
		s:=big.Int{}
		//拆分我们的pubKey,平均分，前半部分给X,后半部分给Y
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])
	//定义两个辅助的big,int
		X:=big.Int{}
		Y:=big.Int{}
	//拆分我们的pubKey,平均分，前半部分给X,后半部分给Y
		X.SetBytes(pubKey[:len(pubKey)/2])
		Y.SetBytes(pubKey[len(pubKey)/2:])

		pubKeyOrigin := ecdsa.PublicKey{elliptic.P256(),&X,&Y}
	//Verify
		if !ecdsa.Verify(&pubKeyOrigin,dataHash,&r,&s){
		return false
		}
	}
	return true
}

func (tx Transacation) String() string {
	var lines []string
	lines =append(lines ,fmt.Sprintf("--- Transaction %x:",tx.TXID))

	for i,input := range tx.TXInputs {
		lines =append(lines,fmt.Sprintf("    Input   %d:",i))
		lines =append(lines,fmt.Sprintf("     TXid     %x",input.TXid))
		lines =append(lines,fmt.Sprintf("      Out     %d",input.Index))
		lines =append(lines,fmt.Sprintf("   Signature  %x",input.Signature))
		lines =append(lines,fmt.Sprintf("     PubKey   %x",input.PubKey))
	}
	for i,output := range tx.TXOutputs {
		lines =append(lines,fmt.Sprintf("              %d:",i))
		lines =append(lines,fmt.Sprintf("     Value      %f",output.Value))
		lines =append(lines,fmt.Sprintf("     Script     %x",output.PubKeyHash))
	}
	return strings.Join(lines,"\n")
}


















//判断是否为挖矿交易

func (tx *Transacation ) IsCoinbase() bool {
	//交易input只有一个
	//交易id为空
	//交易的index 为-1
	if len(tx.TXInputs)==1 && len(tx.TXInputs[0].TXid)==0 && tx.TXInputs[0].Index == -1 {
//		input:=tx.TXInputs[0]
//		if !bytes.Equal(input.TXid,[]byte{}) ||input.Index!=-1{
//			return false
//		}
		return true
	}
	return false
}
