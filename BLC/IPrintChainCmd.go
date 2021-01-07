package BLC

import (
	"fmt"
	"os"
)

func PrintChain_Cmd()  {
	blockChainObject := GetBlockchainObject()
	if blockChainObject == nil{
		fmt.Println("[!]无法获取数据库，请检查！")
		os.Exit(1)
	}
	blockChainObject.PrintChain()
	defer blockChainObject.DB.Close()
}
