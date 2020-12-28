package BLC

type TXInput struct {
	TxID []byte //交易ID
	Vout int //存储TxOutput的Vout里面的索引
	ScriptSiq string //用户名
}

func (txi *TXInput) UnLockWithAddress(add string) bool {
	return txi.ScriptSiq == add
}
