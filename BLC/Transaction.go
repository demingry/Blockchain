package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
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


/*
	创世区块的Transaction
*/
func NewCoinBaseTransaction(addr string) *Transaction {
	txInput := &TXInput{[]byte{},-1,"Genesis block"}
	txOutput := &TXOutput{10,addr}
	txCoinBase := &Transaction{[]byte{},[]*TXInput{txInput},[]*TXOutput{txOutput}}

	txCoinBase.SetTxID()
	return txCoinBase
}



/*
	新区快的Transaction
*/
func NewSimpleTransaction(from,to string,amount int64,bc *Blockchain,txs []*Transaction)*Transaction{
	var txInputs [] *TXInput
	var txOutputs [] *TXOutput

	balance, spendableUTXO := bc.FindSpendableUTXOs(from,amount,txs)


	//代表消费
	for txID,indexArray:=range spendableUTXO{
		txIDBytes,_:=hex.DecodeString(txID)
		for _,index:=range indexArray{
			txInput := &TXInput{txIDBytes,index,from}
			txInputs = append(txInputs,txInput)
		}
	}

	//转账
	txOutput1 := &TXOutput{amount, to}
	txOutputs = append(txOutputs, txOutput1)

	//找零

	txOutput2 := &TXOutput{balance - amount, from}


	txOutputs = append(txOutputs, txOutput2)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	//设置hash值
	tx.SetTxID()
	return tx
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



/*
	判断是否创世区块Transaction
*/
func (tx *Transaction) IsCoinbaseTransaction() bool {

	return len(tx.Vins[0].TxID) == 0 && tx.Vins[0].Vout == -1
}