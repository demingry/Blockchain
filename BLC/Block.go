package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	//区块链高度
	Height int64
	//上一个区块Hash
	PrevBlockHash []byte
	//交易数据，一个transaction代表一次交易，结构体数组表示区块链中多个transaction交易
	TxS []*Transaction
	//时间戳
	Timestamp int64
	//Hash
	Hash []byte
	//工作量证明穷举
	Nonce int64
}

/*
func (b *Block)SetHash()  {
	//Height转化字节数组
	heightBytes := IntToHex(b.Height)

	//Timestamp转化字节数组
	timeString := strconv.FormatInt(b.Timestamp,2)
	timeBytes := []byte(timeString)

	//拼接属性
	blockBytes := bytes.Join([][]byte{heightBytes,b.PrevBlockHash,b.Data,timeBytes,b.Hash},[]byte{})

	//生成Hash
	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:]
}
*/

/*
	创建创世区块
*/
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs,1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}


/*
	新建区块
*/
func NewBlock(txs []*Transaction,height int64,prevblockhash []byte) *Block {
	block := &Block{height,prevblockhash,txs,time.Now().Unix(),nil,0}
	//block.SetHash()

	//创建工作量证明对象
	pow := NewProofofWork(block)
	//穷举直到正确hash被计算出来
	hash,nonce := pow.Run()

	//设置运行出来的hash和穷举数nonce
	block.Hash = hash
	block.Nonce = nonce
	return block
}

/*
	序列化函数
*/
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}


/*
	反序列化函数
*/
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}


/*
	将transaction转化成字节数组
*/
func (b *Block) HashTransaction() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _,tx := range b.TxS {
		txHashes = append(txHashes,tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes,[]byte{}))

	return txHash[:]
}
