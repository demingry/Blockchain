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
	cli := BLC.CLI{}
	cli.Run()
}


/*
	Transaction: TxHash []byte; Vins []*TXInput; Vouts []*TXOutput
	TXInput: TxID []byte; Vout int; ScriptSiq string
	TXOutput: Value int64; ScriptPubKey string
*/
