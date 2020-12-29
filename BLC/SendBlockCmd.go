package BLC

import (
	"fmt"
	"os"
)

func SendBlock_Cmd(from, to, amount [] string) {
	if !DbExists() {
		fmt.Println("数据库不存在。。。")
		os.Exit(1)
	}
	blockchain := GetBlockchainObject()

	blockchain.MineNewBlock(from, to, amount)
	defer blockchain.DB.Close()
}