package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Remove("./blockchain.db")
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "./blockchain.db")
	if err != nil {
		log.Fatal(err)
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
	/*createTableSQL2 := `CREATE TABLE IF NOT EXISTS transactions (
		unspent_hash TEXT PRIMARY KEY,
		block_hash TEXT,
	);`
	if _, err := db.Exec(createTableSQL2); err != nil {
		log.Fatal(err)
	}*/
	/*
		block := NewBlock("Sample Block Data", []byte{})
		insertBlock(db, block)
		fmt.Println("Block inserted into database.")*/
}
