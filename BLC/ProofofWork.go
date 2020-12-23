package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofofWork struct {
	Block *Block  //当前要验证的区块
	target *big.Int //大数据存储
}

const targetBit = 20

//拼接数据字符数组
func (p *ProofofWork)PrepareData(nounce int) []byte {
	data := bytes.Join(
		[][]byte{
			p.Block.PrevBlockHash,
			p.Block.Data,
			IntToHex(int64(p.Block.Height)),
			IntToHex(p.Block.Timestamp),
			IntToHex(int64(targetBit)),
			IntToHex(int64(nounce)),
		},[]byte{},
		)
	return data
}


//验证区块hash方法
func (p *ProofofWork) IsValid() bool {

	var hashInt big.Int
	hashInt.SetBytes(p.Block.Hash)

	if p.target.Cmp(&hashInt)==1 {
		return true
	}
	return false
}

//遍历穷举直到计算hash的值小于target
func (p *ProofofWork)Run() ([]byte,int64) {
	var hashInt big.Int
	var hash [32]byte
	nouce := 0
	for {
		databytes := p.PrepareData(nouce)
		hash = sha256.Sum256(databytes)
		fmt.Printf("\r%x",hash)
		hashInt.SetBytes(hash[:])

		//判断hashInt是否小于target
		if p.target.Cmp(&hashInt)==1 {
			break
		}
		
		nouce = nouce + 1
	}
	fmt.Printf("\n")
	return hash[:],int64(nouce)
}

func NewProofofWork(b *Block) *ProofofWork {
	target := big.NewInt(1)
	target = target.Lsh(target,uint(256-targetBit))
	return &ProofofWork{b,target}
}
