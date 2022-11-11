package main

import "fmt"

func conTinues()  {
	OUTPUT:	
	for  j:=0;j<3;j++{

			for i :=0; i < 5; i++{
				fmt.Println(i)
				continue OUTPUT
			}
	}
}