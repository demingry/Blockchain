package BLC

import (
	"time"
)

type Block struct {
	//区块链高度
	Height int64
	//上一个区块Hash
	PrevBlockHash []byte
	//交易数据
	Data []byte
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

func CreateGenesisBlock(data string) *Block {
	return NewBlock(data,1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}

func NewBlock(data string,height int64,prevblockhash []byte) *Block {
	block := &Block{height,prevblockhash,[]byte(data),time.Now().Unix(),nil,0}
	/*block.SetHash()*/

	//创建工作量证明对象
	pow := NewProofofWork(block)
	//穷举直到正确hash被计算出来
	hash,nonce := pow.Run()

	/*设置运行出来的hash和穷举数nonce*/
	block.Hash = hash
	block.Nonce = nonce
	return block
}
