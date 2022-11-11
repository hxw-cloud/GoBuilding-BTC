package main

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func hash()  {
	file,_:=os.OpenFile("./hash.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	var data ="hello world"
	defer file.Close()
	for i := 0; i < 100; i++ {
		hash:=sha256.Sum256([]byte(data+string(i)))
		fmt.Printf("hash : %x\n",hash[:])
		_,err:=file.Write(hash[:])
		if err!=nil{
			return
		}
	}
}