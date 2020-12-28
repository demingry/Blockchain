package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
	"blockchain/BLC/CLIex"
)


type CLI struct {

}

/*
	打印使用信息
*/
func PrintUsage() {
	fmt.Println("[!]Usage:")
	fmt.Println("\taddblock -addr DATA --交易数据")
	fmt.Println("\tprintchain --输出区块信息")
	fmt.Println("\tcreateblockchain --创建创世区块")
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
	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)

	//设置标签后的参数
	flagAddBlockData := addBlockCmd.String("data","georgedeming","block data")
	flagCreateBlockChainAddr := createBlockChainCmd.String("addr","Genesis block data..","创世区块交易数据")
	flagFromData := sendBlockCmd.String("from","","转账源地址")
	flagToData := sendBlockCmd.String("to","","转账目的地址")
	flagAmountData := sendBlockCmd.String("amount","","转账金额")
	flagGetBalanceData := getBalanceCmd.String("addr","","查询账户余额")

	/*
		解析:CLI命令失误退出
	*/
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
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			os.Exit(1)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if  err != nil {
			os.Exit(1)
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
		CLIex.AddBlock_Cmd()
	}

	if printChainCmd.Parsed() {
		CLIex.PrintChain_Cmd()
	}

	if createBlockChainCmd.Parsed(){
		if *flagCreateBlockChainAddr == ""{
			PrintUsage()
			os.Exit(1)
		}
		CreateBlockchainWithGenesisBlock(*flagCreateBlockChainAddr)
	}

	if sendBlockCmd.Parsed(){
		if *flagFromData == "" || *flagToData =="" ||*flagAmountData == "" {
			PrintUsage()
			os.Exit(1)
		}
		from:=JSONToArray(*flagFromData)
		to:=JSONToArray(*flagToData)
		amount:=JSONToArray(*flagAmountData)

		CLIex.SendBlock_Cmd(from,to,amount)
	}

	if getBalanceCmd.Parsed(){
		if *flagGetBalanceData == ""{
			fmt.Println("查询地址不能为空")
			PrintUsage()
			os.Exit(1)
		}
		CLIex.GetBalance_Cmd(*flagGetBalanceData)
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
