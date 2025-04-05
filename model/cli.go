package model

import (
	"flag"
	"fmt"
	"os"
)

type CLI struct{}

func (cli *CLI) Run() {
	cli.validateArgs()
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockChain := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getBalance := flag.NewFlagSet("getbalance", flag.ExitOnError)
	send := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := send.String("from", "", "Source address")
	sendTo := send.String("to", "", "Destination address")
	sendAmount := send.Int("amount", 0, "Amount to send")
	createBlockChainAddress := createBlockChain.String("address", "", "Coinbase address")
	getBalanceAddress := getBalance.String("address", "", "Wallet address")
	var err error
	switch os.Args[1] {
	case "printchain":
		err = printChainCmd.Parse(os.Args[2:])
	case "createblockchain":
		err = createBlockChain.Parse(os.Args[2:])
	case "getbalance":
		err = getBalance.Parse(os.Args[2:])
	case "send":
		err = send.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("Error parsing command line arguments:", err)
		os.Exit(1)
	}
	if createBlockChain.Parsed() {
		if *createBlockChainAddress == "" {
			createBlockChain.Usage()
			os.Exit(1)
		}
		CreateBlockChain(*createBlockChainAddress)
		fmt.Println("Blockchain created successfully!")
	} else if getBalance.Parsed() {
		if *getBalanceAddress == "" {
			getBalance.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	} else if send.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			send.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
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

func (cli *CLI) getBalance(address string) {
	bc := NewBlockChain(address)
	defer bc.db.Close()
	balance := 0
	utxos := bc.FindUTXO(address)
	for _, txo := range utxos {
		balance += txo.Value
	}
	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockChain(from)
	defer bc.db.Close()
	tx := newUtxoTx(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Transaction successful!")
}
