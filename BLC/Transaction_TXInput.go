package BLC

type TXInput struct {
	TxID []byte
	Vout int
	ScriptSiq string
}

func (txi *TXInput) UnLockWithAddress(add string) bool {
	return txi.ScriptSiq == add
}
