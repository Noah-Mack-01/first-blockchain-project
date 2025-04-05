package model

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

const targetBits = 10

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{block, target}
}

func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.HashTransactions(),
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
}

func (pow *ProofOfWork) Run() (int, []byte) {
	nonce := 0
	var hash [32]byte
	var bigInt big.Int

	for nonce < math.MaxInt64 {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		bigInt.SetBytes(hash[:])

		if bigInt.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hash [32]byte
	var bigInt big.Int

	data := pow.PrepareData(pow.Block.Nonce)
	hash = sha256.Sum256(data)
	bigInt.SetBytes(hash[:])

	return bigInt.Cmp(pow.Target) == -1
}
