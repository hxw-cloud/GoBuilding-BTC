package main


//import (
//	"fmt"
//)

func main()  {
	// 注册模型
	bc:=NewBlockChain("13iR87DSAm3itd3qt9vtvXPdKqs63rQVf8")
	cli:=CLI{bc}
	cli.Run()
}