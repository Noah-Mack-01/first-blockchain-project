package model

import (
	"bytes"
	"encoding/gob"
)

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encode := gob.NewEncoder(&result)
	err := encode.Encode(b)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return &block
}
