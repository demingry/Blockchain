package BLC

import (
	"fmt"
	"os"
)

func GetBalance_Cmd(addr string) {
	fmt.Println("查询余额：",addr)
	bc := GetBlockchainObject()

	if bc == nil{
		fmt.Println("数据库不存在，无法查询。。")
		os.Exit(1)
	}
	defer bc.DB.Close()

	balance:=bc.GetBalance(addr,[]*Transaction{})
	fmt.Printf("%s,一共有%d个Token\n",addr,balance)
}
