package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

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


type UTXO struct {
	TxID   [] byte  //当前Transaction的交易ID
	Index  int      //下标索引
	Output *TXOutput //
}

func NewCoinBaseTransaction(addr string) *Transaction {
	txInput := &TXInput{[]byte{},-1,"Genesis block"}
	txOutput := &TXOutput{10,addr}
	txCoinBase := &Transaction{[]byte{},[]*TXInput{txInput},[]*TXOutput{txOutput}}

	txCoinBase.SetTxID()
	return txCoinBase
}

func (tx *Transaction) SetTxID()  {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	buffBytes := bytes.Join([][]byte{IntToHex(time.Now().Unix()),buff.Bytes()},[]byte{})
	hash := sha256.Sum256(buffBytes)
	tx.TxHash = hash[:]
}

func (tx *Transaction) IsCoinbaseTransaction() bool {

	return false
}