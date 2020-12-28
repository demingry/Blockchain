package CLIex

import (
	"blockchain/BLC"
	"fmt"
	"os"
)

func SendBlock_Cmd(from, to, amount [] string) {
	if !BLC.DbExists() {
		fmt.Println("数据库不存在。。。")
		os.Exit(1)
	}
	blockchain := BLC.GetBlockchainObject()

	blockchain.MineNewBlock(from, to, amount)
	defer blockchain.DB.Close()
}