package BLC

/*
Transaction 创建分两种情况:
1.创世区块创建时的Transaction
2.转账时产生的Transaction
*/

type Transaction struct {
	TxID []byte
	Vins []*TXInput
	Vouts []*TXOutput
}