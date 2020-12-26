package BLC

/*
Transaction 创建分两种情况:
1.创世区块创建时的Transaction
2.转账时产生的Transaction
*/

type Transaction struct {
	TxHash []byte //交易hash
	Vins []*TXInput //输入
	Vouts []*TXOutput //输出
}