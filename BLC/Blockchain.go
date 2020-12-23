package BLC

import (
	"bolt"
	"fmt"
	"log"
	"math/big"
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
			log.Panic(err)
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
	var currentHash = bc.Tip //需要遍历的当前hash值
	var block *Block //需要遍历的当前的区块
	var count = 0
	for {
		err := bc.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blockTableName))

			if b != nil {
				blockBytes := b.Get(currentHash)
				block = DeserializeBlock(blockBytes) //获取响应hash的区块字节数组数据
				fmt.Printf("Current: %d\n",count)
				fmt.Printf("Height: %d\n",block.Height)
				fmt.Printf("PrevBlockHash: %x\n",block.PrevBlockHash)
				fmt.Printf("Data: %s\n",block.Data)
				fmt.Printf("Timestamp: %d\n",block.Timestamp)
				fmt.Printf("Hash: %x\n",block.Hash)
				fmt.Printf("Nonce: %d\n",block.Nonce)
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		//如果此区块的前一个区块的hash为0，则是创世区块
		if big.NewInt(0).Cmp(&hashInt)==0 {
			break
		}
		currentHash = block.PrevBlockHash //更新需要遍历的区块hash值
		count++
	}
}
