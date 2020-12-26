package BLC

import (
	"bolt"
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
func CreateBlockchainWithGenesisBlock(tsx []*Transaction) {

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

			genesisBlock := CreateGenesisBlock(tsx)
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
		fmt.Printf("Transaction%v\n",block.TxS)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)

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
