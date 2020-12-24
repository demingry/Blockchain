package BLC

import (
	"bolt"
	"log"
)

type BlockchainIterator struct {
	CurrentHash []byte
	DB *bolt.DB
}

/*
	创建区块链迭代对象
*/
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.Tip,bc.DB}
}

/*
	获取区块数据
*/
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


