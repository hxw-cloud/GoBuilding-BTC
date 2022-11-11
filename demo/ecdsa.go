package  main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
)

//演示如何使用ecdic生成公钥和私钥
//校验签名

func ecdic(){
	//创建曲线
	curve :=elliptic.P256()
	//生成私钥
	privateKey,err:=ecdsa.GenerateKey(curve,rand.Reader)
	if err!=nil{
		log.Panic("")
	}
	//生成公钥
	pubKey :=privateKey.PublicKey
	data :="hello world!"
	hash :=sha256.Sum256([]byte(data))
	fmt.Printf("pubKey : %v\n",pubKey)

	//签名
	r,s,err:=ecdsa.Sign(rand.Reader,privateKey,hash[:])
	if err!=nil{
		log.Panic("")
	}


	fmt.Printf("r : %v , len : %d\n",r.Bytes(),len(r.Bytes()))
	fmt.Printf("s : %v , len : %d\n",s.Bytes(),len(s.Bytes()))
	//把r,s进行序列化传输
	signature :=append(r.Bytes(),s.Bytes()...)
	fmt.Printf("signature : %v\n",signature)

	//定义两个辅助的big,int
	s1:=big.Int{}
	r1:=big.Int{}
	//拆分我们的signature,平均分，前半部分给r,后半部分给s
	r1.SetBytes(signature[:len(signature)/2])
	s1.SetBytes(signature[len(signature)/2:])

	//校验需要三个东西： 数据，签名，公钥
	res:=ecdsa.Verify(&pubKey,hash[:],&r1,&s1)
	fmt.Println("校验结果 : ",res)
}