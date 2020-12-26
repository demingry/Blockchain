package BLC

type TXOutput struct {
	Value        int64
	ScriptPubKey string
}

func (txo *TXOutput) UnLockWithAddress(add string) bool {
	return txo.ScriptPubKey == add
}
