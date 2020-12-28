package CLIex

import (
	"fmt"
	"os"
	"blockchain/BLC"
)

func AddBlock_Cmd() {
	blockChainObject := BLC.GetBlockchainObject()
	if blockChainObject == nil{
		fmt.Println("[!]无法获取数据库，请检查！")
		os.Exit(1)
	}
	blockChainObject.AddBlockToBlockChain([]*BLC.Transaction{})
	defer blockChainObject.DB.Close()
}
