package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

//使用gob实现序列化（编码）
//使用gob进行反序列化

type Person struct{
	Name string
	Age uint64
}

func gobDemo()  {
	var xiaoMing Person
	xiaoMing.Name="小明"
	xiaoMing.Age=120
	//编码数据放到buffer
	fmt.Println(xiaoMing)
	var buffer bytes.Buffer
	//使用gob进行序列化编码
	//定义一个编码器
	//使用编码器进行编码
	encoder :=gob.NewEncoder(&buffer)
	err:=encoder.Encode(&xiaoMing)
	if err!=nil{
		log.Panic("编码出错，小明不知去向")
	}
	fmt.Printf("编码后的小明： %v\n",buffer.Bytes())
	//使用gob进行反序列化（解码）得到Person结构
	//定义一个解码器
	//使用解码器解码
	decoder :=gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	var daMing Person
	err=decoder.Decode(&daMing)
	if err!=nil{
		log.Panic("解码出错，小明不知去向")
	}

	fmt.Printf("编码后的小明： %v\n",daMing)
}