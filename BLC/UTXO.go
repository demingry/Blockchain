package BLC

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
)



type UTXO struct {
	TxID   [] byte  //当前Transaction的交易ID
	Index  int      //下标索引
	Output *TXOutput //
}

//找到所有未花费的交易输出
func (bc *Blockchain) UnUTXOs(addr string, txs []*Transaction) []*UTXO {
	/*
	   1.先遍历未打包的交易(参数txs)，找出未花费的Output。
	   2.遍历数据库，获取每个块中的Transaction，找出未花费的Output。
	*/
	var unUTXOs []*UTXO                      //未花费
	spentTxOutputs := make(map[string][]int) //存储已经花费

	//1.添加先从txs遍历，查找未花费
	//for i, tx := range txs {
	for i:=len(txs)-1;i>=0;i--{
		unUTXOs = caculate(txs[i], addr, spentTxOutputs, unUTXOs)
	}

	bcIterator := bc.Iterator()
	for {
		block := bcIterator.Next()
		//统计未花费
		//2.获取block中的每个Transaction
		for i := len(block.TxS) - 1; i >= 0; i-- {
			unUTXOs = caculate(block.TxS[i], addr, spentTxOutputs, unUTXOs)
		}

		//结束迭代
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}
	}
	return unUTXOs
}


func (bc *Blockchain) GetBalance(address string, txs []*Transaction) int64 {
	unUTXOs := bc.UnUTXOs(address, txs)
	var amount int64
	for _, utxo := range unUTXOs {
		amount = amount + utxo.Output.Value
	}
	return amount

}


//转账时查获在可用的UTXO
func (bc *Blockchain) FindSpendableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int) {
	/*
	   1.获取所有的UTXO
	   2.遍历UTXO

	   返回值：map[hash]{index}
	*/

	var balance int64
	utxos := bc.UnUTXOs(from, txs)
	spendableUTXO := make(map[string][]int)
	for _, utxo := range utxos {
		balance += utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxID)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if balance >= amount {
			break
		}
	}
	if balance < amount {
		fmt.Printf("%s 余额不足。。总额：%d，需要：%d\n", from,balance,amount)
		os.Exit(1)
	}
	return balance, spendableUTXO
}


func caculate(tx *Transaction, addr string, spentTxOutputs map[string][]int, unUTXOs []*UTXO) []*UTXO {
	//2.先遍历TxInputs，表示花费
	if !tx.IsCoinbaseTransaction() {
		for _, in := range tx.Vins {
			//如果解锁
			if in.UnLockWithAddress(addr) {
				key := hex.EncodeToString(in.TxID)
				spentTxOutputs[key] = append(spentTxOutputs[key], in.Vout)
			}
		}
	}

outputs:
	for index, out := range tx.Vouts {
		if out.UnLockWithAddress(addr) {
			if len(spentTxOutputs) != 0 {
				var isSpentUTXO bool

				for txID, indexArray := range spentTxOutputs {
					for _, i := range indexArray {
						if i == index && txID == hex.EncodeToString(tx.TxHash) {
							isSpentUTXO = true
							continue outputs
						}
					}
				}
				if !isSpentUTXO {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				}

			} else {
				utxo := &UTXO{tx.TxHash, index, out}
				unUTXOs = append(unUTXOs, utxo)
			}
		}
	}
	return unUTXOs
}
