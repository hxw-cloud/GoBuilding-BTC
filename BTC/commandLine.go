package main

import (
	"fmt"
)
//func (cli *CLI)PrintBlockChain()  {
//	cli.bc.PrintChain()
//	fmt.Print("打印区块链完成\n")
//}

func (cli *CLI)PrintBlockChainReverse()  {
	bc:=cli.bc
	it :=bc.NewIterator()
	for {
		block:=it.Next()
		for _,tx := range block.Transacations{
			fmt.Println(tx)
		}
		/*fmt.Println("=========当前高度为","=========")
		fmt.Printf("版本号： %d\n",block.Version)
		fmt.Printf("前区块Hash值： %x\n",block.PrevHash)
		fmt.Printf("梅克尔根： %x\n",block.MerkelRoot)
		fmt.Printf("当前区块Hash值：%x\n",block.Hash)
		timeFormat :=time.Unix(int64(block.TimeStamp),0).Format("2006-01-02 03:04:05 PM")
		fmt.Printf("时间搓： %s\n",timeFormat)
		fmt.Printf("难度值（未定义，随便写的）： %d\n",block.Difficulty)
		fmt.Printf("随机数： %d\n",block.Nonce)
		fmt.Printf("区块数据：%s\n",block.Transacations[0].TXInputs[0].PubKey)*/
		if len(block.PrevHash)==0{
			fmt.Print("打印结束")
			break
		}
	}


}

//func (cli *CLI ) addBlock(data string)  {
////	cli.bc.AddBlock()
//	fmt.Print("success!")
//}

func (cli *CLI)GetBalance(address string)  {
	//校验地址
	if !IsValiAddress(address){
		fmt.Printf("地址无效 ： %s\n",address)
		return
	}
	//生成公钥hash
	pubKeyHash := GetPubKeyFromAddress(address)
	utxos:=cli.bc.FindUTXOs(pubKeyHash)
	total :=0.0
	for _,utxo:=range utxos{
		total+=utxo.Value
	}
	fmt.Println("余额为：",total)
}

func (cli *CLI) Send(form,to string,amount float64,miner,data string)  {
	//校验地址
	if !IsValiAddress(form){
		fmt.Printf("from 地址无效 ： %s\n",form)
		return
	}
	//校验地址
	if !IsValiAddress(to){
		fmt.Printf("to 地址无效 ： %s\n",to)
		return
	}
	//校验地址
	if !IsValiAddress(miner){
		fmt.Printf("miner 地址无效 ： %s\n",miner)
		return
	}

	//创建挖矿交易
	coinbase :=NewCoinbaseTX(miner,data)
//	fmt.Println("coinbase value",coinbase)
	//创建普通交易
	tx :=NewTransaction(form,to,amount,cli.bc)
	if tx==nil{
		return
	}
	//添加到区块
	cli.bc.AddBlock([]*Transacation{coinbase,tx})
	fmt.Printf("转账成功")
}

func (cli *CLI) NewWallet()  {
	wallet :=NewWallets()
	address := wallet.CreatWallet()
//	address:= wallet.NewAddress()
//	ws := NewWallets()
//
//	for address := range ws.WalletsMap{
		fmt.Printf("address ： %s\n",address)
//	}
}

func (cli  *CLI) ListAddresses()  {
	ws:=NewWallets()
	addresses := ws.ListAllAddresses()
	for _,address := range addresses {
		fmt.Printf("地址: %s\n",address)
	}
}