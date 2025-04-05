package main

import (
	model "noerkrieg/blockchain-tutorial/model"
)

/*func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()
	for {
		block := bci.Next()
		if block == nil {
			break
		}
		fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
		pow := model.NewProofOfWork(block)
		fmt.Printf("POW: %v\n", pow.Validate())
		fmt.Println()
	}
}*/

func main() {
	cli := model.CLI{}
	cli.Run()
}
