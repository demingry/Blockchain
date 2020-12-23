/*
	Author: Wen rui, Zhang
	Duration: 2020-2021
	E-Mail: georgedeming@outlook.com
*/

package main

import (
	"blockchain/BLC"
	"fmt"
)

func main()  {
	Blockchain := BLC.CreateBlockchainWithGenesisBlock()
	Blockchain.AddBlockToBlockChain("send $100 to Zhang")
	Blockchain.AddBlockToBlockChain("send $100 to Ge")
	defer Blockchain.DB.Close()
	Blockchain.PrintChain()
	fmt.Scanf("this")
}
