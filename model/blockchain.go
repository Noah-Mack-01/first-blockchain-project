package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type BlockChain struct {
	tip []byte
	db  *sql.DB
}

type BlockChainIterator struct {
	currentHash []byte
	db          *sql.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.tip, bc.db}
}

func (bc *BlockChain) Close() {
	bc.db.Close()
}

func (bci *BlockChainIterator) Next() *Block {
	var blockData []byte
	err := bci.db.QueryRow("SELECT serialized FROM blocks WHERE hash = ?", bci.currentHash).Scan(&blockData)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	} else if err == sql.ErrNoRows {
		return nil
	}
	block := Deserialize(blockData)
	if block == nil {
		return nil
	}
	bci.currentHash = block.PrevBlockHash
	return block
}

func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte
	// query for all columns in table blocks
	err := bc.db.QueryRow("SELECT hash FROM blocks ORDER BY timestamp DESC LIMIT 1").Scan(&lastHash)
	if err != nil {
		panic(err)
	}
	block := NewBlock(data, lastHash)
	res, err := bc.db.Exec("INSERT INTO blocks (hash, serialized, timestamp) VALUES (?, ?, ?)", block.Hash, block.Serialize(), block.Timestamp)
	if err != nil {
		panic(err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}
	bc.tip = block.Hash
	log.Printf("Added block with hash: %d", block.Hash)
}

func NewBlockChain() *BlockChain {
	var tip []byte
	db, err := sql.Open("sqlite3", "./blockchain.db")
	if err != nil {
		panic(err)
	}

	var hash []byte
	err = db.QueryRow("SELECT hash FROM blocks ORDER BY timestamp DESC LIMIT 1").Scan(&hash)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	if hash == nil {
		genesis := GenesisBlock()
		_, err := db.Exec(
			"INSERT INTO blocks (hash, serialized, timestamp) VALUES (?, ?, ?)",
			genesis.Hash, genesis.Serialize(), genesis.Timestamp)
		if err != nil {
			panic(err)
		}
		tip = genesis.Hash
	} else {
		tip = hash
	}

	return &BlockChain{tip, db}
}
