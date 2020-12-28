package BLC

type TXOutput struct {
	Value        int64 //金额
	ScriptPubKey string //用户名
}

func (txo *TXOutput) UnLockWithAddress(add string) bool {
	return txo.ScriptPubKey == add
}
