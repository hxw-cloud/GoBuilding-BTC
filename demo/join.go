package main

import (
	"bytes"
	"fmt"
	"strings"
)

func Join()  {
	str1 :=[]string{"hello world","huxinwei","!"}
	res:=strings.Join(str1,"")
	fmt.Println(res)
	res1:=bytes.Join([][]byte{[]byte("hello"),[]byte("world")},[]byte("+"))
	fmt.Printf("%s\n",res1)
}