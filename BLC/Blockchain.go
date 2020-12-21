package BLC

import (
	"bolt"
	"log"
)

const (
	dbName = "blockedchain.db" //数据库名
	blockTableName = "blocks" //数据表名
)

type Blockchain struct {
	Tip []byte
	DB *bolt.DB
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



//创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain {

	db,err := bolt.Open(dbName,0600,nil)
	if err != nil {
		log.Fatal(err)
	}

	var blockhash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b,err := tx.CreateBucket([]byte(blockTableName)) //创建数据表
		if err != nil {
			panic(err)
		}

		//数据表不存在时,创建创世区块
		if b == nil {

			genesisBlock := CreateGenesisBlock("genesis block creating...")
			//将创世区块存储到表中
			err := b.Put(genesisBlock.Hash,genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			blockhash = genesisBlock.Hash

			//存储最新的区块的hash
			err = b.Put([]byte("l"),genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	return &Blockchain{blockhash,db}
}
