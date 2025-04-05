package model

import (
	"database/sql"
	"encoding/hex"
	"log"
	"os"

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

func (bc *BlockChain) MineBlock(transactions []*Transaction) {
	var lastHash []byte
	// query for all columns in table blocks
	err := bc.db.QueryRow("SELECT hash FROM blocks ORDER BY timestamp DESC LIMIT 1").Scan(&lastHash)
	if err != nil {
		panic(err)
	}
	block := NewBlock(transactions, lastHash)
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

func NewBlockChain(address string) *BlockChain {
	var tip []byte
	if info, _ := os.Stat("./blockchain.db"); info == nil {
		log.Panic("No blockchain found, create one.")
	}
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
		cbtx := NewCoinbaseTX(address, "Genesis Block")
		genesis := GenesisBlock(cbtx)
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

func CreateBlockChain(address string) {
	if info, _ := os.Stat("./blockchain.db"); info != nil {
		panic("Blockchain already exists")
	}

	db, err := sql.Open("sqlite3", "./blockchain.db")
	// creating new DB

	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create a table for blocks if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS blocks (
		hash TEXT PRIMARY KEY,
		serialized TEXT,
		timestamp INTEGER
);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal(err)
	}

	cbtx := NewCoinbaseTX(address, "Genesis Block")
	genesis := GenesisBlock(cbtx)

	res, err := db.Exec(
		"INSERT INTO blocks (hash, serialized, timestamp) VALUES (?, ?, ?)",
		genesis.Hash, genesis.Serialize(), genesis.Timestamp)
	if err != nil {
		panic(err)
	}

	blockID, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	log.Printf("Genesis block created with ID: %d", blockID)
}

func (bc *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspent []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for block := bci.Next(); block != nil; block = bci.Next() {
		for _, tx := range block.Transactions {
			TXID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[TXID] != nil {
					for _, spentOutIdx := range spentTXOs[TXID] {
						if outIdx == spentOutIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlockedWith(address) {
					unspent = append(unspent, *tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Id)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspent
}

func (bc *BlockChain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspent := bc.FindUnspentTransactions(address)
	for _, tx := range unspent {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	acc := 0
	for _, tx := range unspentTxs {
		txId := hex.EncodeToString(tx.ID)
		for outId, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && acc < amount {
				acc += out.Value
				unspentOutputs[txId] = append(unspentOutputs[txId], outId)
			}
			if acc >= amount {
				break
			}
		}
	}
	return acc, unspentOutputs
}
