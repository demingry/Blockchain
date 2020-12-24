/*
	Author: Wen rui, Zhang
	Duration: 2020-2021
	E-Mail: georgedeming@outlook.com
*/

package main

import (
	"blockchain/BLC"
)

func main()  {
	Blockchain := BLC.CreateBlockchainWithGenesisBlock()
	Blockchain.AddBlockToBlockChain(BLC.CreateRandomString(6))
	Blockchain.PrintChain()
	defer Blockchain.DB.Close()
}
