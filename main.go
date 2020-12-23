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
	for {
		Blockchain.AddBlockToBlockChain(BLC.CreateRandomString(6))
	}
	defer Blockchain.DB.Close()
	Blockchain.PrintChain()
	fmt.Scanf("this")
}
