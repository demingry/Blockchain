package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)


type CLI struct {

}

/*
	打印使用信息
*/
func PrintUsage() {
	fmt.Println("[!]Usage:")
	fmt.Println("\taddblock -data DATA --交易数据")
	fmt.Println("\tprintchain --输出区块信息")
}

/*
	主运行方法
*/
func (cli *CLI) Run()  {

	isValidArgs()

	//创建flagset标签对象
	addBlockCmd := flag.NewFlagSet("addblock",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockChainCmd:=flag.NewFlagSet("createblockchain",flag.ExitOnError)

	//设置标签后的参数
	flagAddBlockData := addBlockCmd.String("data","georgedeming","block data")
	flagCreateBlockChainData := createBlockChainCmd.String("data","Genesis block data..","创世区块交易数据")

	switch os.Args[1] {
	case "addblock": //添加区块
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain": //打印区块
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain": //创建创世区块
		err :=createBlockChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	default:
		PrintUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == ""{ //-data 后的数据为空
			PrintUsage()
			os.Exit(1)
		}
		blockChainObject := GetBlockchainObject()
		if blockChainObject == nil{
			fmt.Println("[!]无法获取数据库，请检查！")
			os.Exit(1)
		}
		blockChainObject.AddBlockToBlockChain(CreateRandomString(6))
		defer blockChainObject.DB.Close()
	}

	if printChainCmd.Parsed() {
		blockChainObject := GetBlockchainObject()
		if blockChainObject == nil{
			fmt.Println("[!]无法获取数据库，请检查！")
			os.Exit(1)
		}
		blockChainObject.PrintChain()
		defer blockChainObject.DB.Close()
	}

	if createBlockChainCmd.Parsed(){
		if *flagCreateBlockChainData == ""{
			PrintUsage()
			os.Exit(1)
		}
		CreateBlockchainWithGenesisBlock()
	}
}

/*
	验证参数,为避免数组越界，无操作参数将退出
*/
func isValidArgs(){
	if len(os.Args) < 2{
		PrintUsage()
		os.Exit(1)
	}
}
