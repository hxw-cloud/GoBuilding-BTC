package main

import (
	"encoding/base64"
	"fmt"
)

//base64 生成地址

func address()  {
	//演示base64编码
	input :=[]byte("生成地址")
	encodeString:=base64.StdEncoding.EncodeToString(input)
	fmt.Println(encodeString)
	//对上面的编码进行base64编码
	decodeBytes,_:=base64.StdEncoding.DecodeString(encodeString)
	fmt.Println(decodeBytes)
	fmt.Println()
	//如过要用在url中，需要用URLEncoding
	uENC:=base64.URLEncoding.EncodeToString([]byte(input))
	fmt.Println(uENC)
	uDec,_:=base64.URLEncoding.DecodeString(uENC)
	fmt.Println(uDec)
}