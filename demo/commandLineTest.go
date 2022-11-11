package main

import (
	"fmt"
	"os"
)

func commandLineTest()  {
	len1 :=len(os.Args)
	fmt.Print("cmd len :",len1,"\n")
	for i ,cmd:= range os.Args{
		fmt.Println("第",i,"数字",cmd)
	}
}