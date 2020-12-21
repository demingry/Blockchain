package BLC

type Blockchain struct {
	Blocks []*Block
}

//将新区快增加到区块链中
func (blc *Blockchain)AddNewBlockToChain(data string,height int64,prevhash []byte)  {
	//创建新区块
	newblock := NewBlock(data,height,prevhash)
	//将区块添加
	blc.Blocks = append(blc.Blocks,newblock)
}


//创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain {
	genesisBlock := CreateGenesisBlock("Genesis Block Created ......")
	return &Blockchain{[]*Block{genesisBlock}}
}
