package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

// subsidy is the amount of coins a miner gets for mining a block.
const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// Output stores coins.
// Storing means locking.
// We lock
type TXOutput struct {
	Value      int
	PubKeyHash string // user wallet address
}
type TXInput struct {
	Id        []byte // ID of the transaction
	Vout      int    // index of the output
	ScriptSig string // Signature
	// Signature []byte
}

// when miner mines a block, it adds a coinbase transaction.
// coinbase tx is special, doesnt require an output.
// creates one from nowhere.
// genesis block builds first output.

func NewCoinbaseTX(address, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", address)
	}
	txin := TXInput{
		Id:   []byte{},
		Vout: -1,
		// Signature: []byte{},
		ScriptSig: data,
	}
	txout := TXOutput{
		Value:      subsidy,
		PubKeyHash: address,
	}
	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{txin},
		Vout: []TXOutput{txout},
	}
	tx.SetID()
	return &tx
}

func newUtxoTx(from, to string, amount int, bc *BlockChain) *Transaction {
	var Inputs []TXInput
	var Outputs []TXOutput
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	for txid, outs := range validOutputs {
		txID := []byte(txid)
		for _, out := range outs {
			Inputs = append(Inputs, TXInput{txID, out, from})
		}
	}
	Outputs = append(Outputs, TXOutput{amount, to})
	if acc > amount {
		Outputs = append(Outputs, TXOutput{acc - amount, from})
	}
	tx := Transaction{
		ID:   nil,
		Vin:  Inputs,
		Vout: Outputs,
	}
	tx.SetID()
	return &tx

}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.PubKeyHash == unlockingData
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Vout == -1 && len(tx.Vin[0].Id) == 0
}
