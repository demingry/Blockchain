package BLC

import (
	"bolt"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

const (
	dbName         = "blockedchain.db" //数据库名
	blockTableName = "blocks"          //数据表名
)

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

/*
//将新区快增加到区块链中
func (blc *Blockchain)AddNewBlockToChain(data string,height int64,prevhash []byte)  {
	//创建新区块
	newblock := NewBlock(data,height,prevhash)
	//将区块添加
	blc.Blocks = append(blc.Blocks,newblock)
}
*/

/*
	创建带有创世区块的区块链
*/
func CreateBlockchainWithGenesisBlock(addr string) {

	//数据库是否存在
	if dbExists() {
		fmt.Println("数据库已存在")
		return
	}

	fmt.Println("[!]数据库不存在,将创建创世区块")
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucket([]byte(blockTableName)) //创建数据表
		if err != nil {
			log.Panic(err)
		}

		//数据表被新建,创建创世区块
		if b != nil {
			txCoinBase := NewCoinBaseTransaction(addr)
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinBase})
			//将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//存储最新的区块的hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}

/*
	增加新区块到区块链
*/
func (bc *Blockchain) AddBlockToBlockChain(txs []*Transaction) {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		//获取数据表
		b := tx.Bucket([]byte(blockTableName))

		//数据表已存在
		if b != nil {
			blockBytes := b.Get(bc.Tip)                                                //获取最新区块的hash,l:hash
			lastblock := DeserializeBlock(blockBytes)                                  //获取最新区块的反序列化信息
			newBlock := NewBlock(txs, lastblock.Height+1, lastblock.Hash) //新建区块
			err := b.Put(newBlock.Hash, newBlock.Serialize())                          //把最新的区块通过序列化存储到表中
			if err != nil {
				log.Panic(err)
			}
			_ = b.Put([]byte("l"), newBlock.Hash) //把头索引值的hash改为新建的块hash
			bc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

/*
	打印区块链
*/
func (bc *Blockchain) PrintChain() {

	blockchainIterator := bc.Iterator()
	for {
		block := blockchainIterator.Next()

		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("------------------------------------------\n")
		fmt.Printf("Transaction:\n")
		for _,txs := range block.TxS {
			fmt.Printf("TxHash:%x\n",txs.TxHash)
			fmt.Printf("Vins:\n")
			for _,in := range txs.Vins {
				fmt.Printf("%x\n",in.TxID)
				fmt.Printf("%d\n",in.Vout)
				fmt.Printf("%s\n",in.ScriptSiq)
			}
			fmt.Printf("Vouts:\n")
			for _,out := range txs.Vouts {
				fmt.Printf("%d\n",out.Value)
				fmt.Printf("%s\n",out.ScriptPubKey)
			}
		}

		var hahsInt big.Int
		hahsInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hahsInt) == 0 {
			break
		}
	}
}

/*
	判断数据库是否存在
*/
func dbExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

/*
	获取区块对象
*/
func GetBlockchainObject() *Blockchain {

	if !dbExists() {
		fmt.Println("数据库不存在，无法获取区块链")
		return nil
	}

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var blockchain *Blockchain
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("l"))
			blockchain = &Blockchain{hash, db}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return blockchain
}



//找到所有未花费的交易输出
func (bc *Blockchain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	/*
	   1.先遍历未打包的交易(参数txs)，找出未花费的Output。
	   2.遍历数据库，获取每个块中的Transaction，找出未花费的Output。
	*/
	var unUTXOs []*UTXO                      //未花费
	spentTxOutputs := make(map[string][]int) //存储已经花费

	//1.添加先从txs遍历，查找未花费
	//for i, tx := range txs {
	for i:=len(txs)-1;i>=0;i--{
		unUTXOs = caculate(txs[i], address, spentTxOutputs, unUTXOs)
	}

	bcIterator := bc.Iterator()
	for {
		block := bcIterator.Next()
		//统计未花费
		//2.获取block中的每个Transaction
		for i := len(block.TxS) - 1; i >= 0; i-- {
			unUTXOs = caculate(block.TxS[i], address, spentTxOutputs, unUTXOs)
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
	//fmt.Println(from,utxos)
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


func caculate(tx *Transaction, address string, spentTxOutputs map[string][]int, unUTXOs []*UTXO) []*UTXO {
	//2.先遍历TxInputs，表示花费
	if !tx.IsCoinbaseTransaction() {
		for _, in := range tx.Vins {
			//如果解锁
			if in.UnLockWithAddress(address) {
				key := hex.EncodeToString(in.TxID)
				spentTxOutputs[key] = append(spentTxOutputs[key], in.Vout)
			}
		}
	}

outputs:
	for index, out := range tx.Vouts {
		if out.UnLockWithAddress(address) {
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
