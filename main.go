package main

import (
	"flag"
	"fmt"
	model "noerkrieg/blockchain-tutorial/model"
	"os"
)

type CLI struct {
	bc *model.BlockChain
}

func (cli *CLI) Run() {
	cli.validateArgs()
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")
	var err error
	switch os.Args[1] {
	case "addblock":
		err = addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		err = printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("Error parsing command line arguments:", err)
		os.Exit(1)
	}
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Success!!")
}
func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()
	for {
		block := bci.Next()
		if block == nil {
			break
		}
		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
		pow := model.NewProofOfWork(block)
		fmt.Printf("POW: %v\n", pow.Validate())
		fmt.Println()
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func main() {
	bc := model.NewBlockChain()
	defer bc.Close()

	cli := CLI{bc}
	cli.Run()

}
