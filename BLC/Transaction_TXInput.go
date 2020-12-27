package BLC

type TXInput struct {
	TxID []byte //交易ID
	Vout int
	ScriptSiq string
}

func (txi *TXInput) UnLockWithAddress(add string) bool {
	return txi.ScriptSiq == add
}
