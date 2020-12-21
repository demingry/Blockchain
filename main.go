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
	Blockchain.AddNewBlockToChain("Send 100RMB to Zhang",Blockchain.Blocks[len(Blockchain.Blocks)-1].Height + 1,Blockchain.Blocks[len(Blockchain.Blocks)-1].Hash)
	Blockchain.AddNewBlockToChain("Send 200RMB to Zhang",Blockchain.Blocks[len(Blockchain.Blocks)-1].Height + 1,Blockchain.Blocks[len(Blockchain.Blocks)-1].Hash)
	fmt.Println(Blockchain)
}
