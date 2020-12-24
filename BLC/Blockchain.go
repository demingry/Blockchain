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
	dbName = "blockedchain.db" //数据库名
	blockTableName = "blocks" //数据表名
)

type Blockchain struct {
	Tip []byte
	DB *bolt.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	DB *bolt.DB
}

//创建区块链迭代对象
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.Tip,bc.DB}
}

//获取区块数据
func (bci *BlockchainIterator) Next() *Block {
	var block *Block
	err := bci.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b!= nil{
			currentBlockbytes := b.Get(bci.CurrentHash)
			block = DeserializeBlock(currentBlockbytes)
			bci.CurrentHash = block.PrevBlockHash
		}
		return nil
	})
	if err != nil{
		log.Panic(err)
	}
	return block
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

	//数据库是否存在
	if dbExists() {
		fmt.Println("[!]数据库已存在")
		db, err := bolt.Open(dbName, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		//defer db.Close()
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

	fmt.Println("[!]数据库不存在,将创建创世区块")
	db,err := bolt.Open(dbName,0600,nil)
	if err != nil {
		log.Fatal(err)
	}

	var blockhash []byte

	err = db.Update(func(tx *bolt.Tx) error {

		b,err := tx.CreateBucket([]byte(blockTableName)) //创建数据表
		if err != nil {
			log.Panic(err)
		}

		//数据表被新建,创建创世区块
		if b != nil {

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

func (bc *Blockchain)AddBlockToBlockChain(data string)  {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		//获取数据表
		b := tx.Bucket([]byte(blockTableName))

		//数据表已存在
		if b!= nil{
			blockBytes := b.Get(bc.Tip) //获取最新区块的hash,l:hash
			lastblock := DeserializeBlock(blockBytes) //获取最新区块的反序列化信息
			newBlock := NewBlock(data,lastblock.Height+1,lastblock.Hash) //新建区块
			err := b.Put(newBlock.Hash,newBlock.Serialize()) //把最新的区块通过序列化存储到表中
			if err != nil {
				log.Panic(err)
			}
			b.Put([]byte("l"),newBlock.Hash) //把头索引值的hash改为新建的块hash
			bc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (bc *Blockchain)PrintChain() {

	blockchainIterator := bc.Iterator()
	for  {
		block := blockchainIterator.Next()

		fmt.Printf("Height: %d\n",block.Height)
		fmt.Printf("PrevBlockHash: %x\n",block.PrevBlockHash)
		fmt.Printf("Data: %s\n",block.Data)
		fmt.Printf("Timestamp: %s\n",time.Unix(block.Timestamp,0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash: %x\n",block.Hash)
		fmt.Printf("Nonce: %d\n",block.Nonce)

		var hahsInt big.Int
		hahsInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hahsInt) ==0 {
			break
		}
	}
}

func dbExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}
