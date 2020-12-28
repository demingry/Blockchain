package CLIex

import (
	"fmt"
	"os"
	"blockchain/BLC"
)

func GetBalance_Cmd(addr string) {
	fmt.Println("查询余额：",addr)
	bc := BLC.GetBlockchainObject()

	if bc == nil{
		fmt.Println("数据库不存在，无法查询。。")
		os.Exit(1)
	}
	defer bc.DB.Close()

	balance:=bc.GetBalance(addr,[]*BLC.Transaction{})
	fmt.Printf("%s,一共有%d个Token\n",addr,balance)
}
