package main

import (
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *BlockChain
}
const Usage =`
   printChainR	 "print all blockchain data（打印区块链）"
   getBalance --address ADDRESS     "获取指定地址的余额"
   send FROM TO AMOUNT MINER DATA   "由FROM转AMOUNT给TO,由MINER挖矿，同时写入DATA"
   newWallet      "创建一个新的钱包（私钥公钥对）"
   listAddresses   "列举所有的钱包地址"
`
//接收参数的动作   printChain	 "print all blockchain data（正向打印区块链）"

func (cli *CLI) Run()  {
	//得到的所有的参数，命令
	//  addBlock --data "hello world"
	// printChain
	args :=os.Args
	if len(args)<2{
		fmt.Print(Usage)
		return
	}
	//分析命令
	cmd:=args[1]
	switch cmd {
//	case "addBlock":
//		fmt.Println("添加区块")
//		if len(args) ==4 && args[2] == "--data"{
//			//获取命令行数据
//			data := args[3]
//			cli.addBlock(data)
//		}else {
//			fmt.Print("参数使用不当")
//		}
//		//执行相应动作，添加区块
//		//获取数据
//	case "printChain":
//		fmt.Println("正向打印区块")
//		cli.PrintBlockChain()
		//打印区块
	case "printChainR":
		fmt.Println("打印区块")
		cli.PrintBlockChainReverse()
	case "getBalance":
		//打印出余额、地址
		fmt.Println("获取余额")
		if len(args) ==4 && args[2] == "--address"{
			address:=args[3]
			cli.GetBalance(address)
		}
	case "send":
		fmt.Printf("转账开始。。。\n")
		if len(args)!=7{
			fmt.Printf("参数个数错误，请检查！\n")
			fmt.Printf(Usage)
			return
		}else {
			form,to,miner,data:=args[2],args[3],args[5],args[6]
			amount,_:=strconv.ParseFloat(args[4],64)
			cli.Send(form,to,amount,miner,data)
		}
	case "newWallet":
		fmt.Printf("创建新的钱包。。。。。\n")
		cli.NewWallet()
	case "listAddresses":
		fmt.Printf("列举所有返回地址。。。\n")
		cli.ListAddresses()
	default:
		fmt.Println("无效请求请检查")
		fmt.Printf(Usage)
	}
	//执行相应的动作
}