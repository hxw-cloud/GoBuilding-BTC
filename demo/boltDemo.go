package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)
func boltDemo() {
	db,err:=bolt.Open("test.db",0600,nil)
	if err!=nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	_=db.Update(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte("b1"))
		if bucket!=nil{
			//写入数据

		}else {
			//创建
			bucket,err=tx.CreateBucket([]byte("b1"))
			if err!=nil{
				log.Panic("创建数据库失败")
			}
		}
		_=bucket.Put([]byte("1111"),[]byte("hello"))
		_=bucket.Put([]byte("2222"),[]byte("world"))
		return nil
	})
}
