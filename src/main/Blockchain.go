package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const blocksBucket = "ismailchain.bucket"
const dbFile = "ismailchain.db"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash
		return nil
	})
}

func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain Already exists")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinBaseTX(address, genesisCoinbaseData)
		genisis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}
		err = b.Put(genisis.Hash, genisis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		tip = genisis.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}
	return &bc
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func NewBlockchain(address string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			cbTx := NewCoinBaseTX(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbTx)
			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
			if err != nil {
			}

		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
	}

	bc := &Blockchain{tip, db}
	return bc
}


func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}

				if tx.IsCoinbase() == false {
					for _, in := range tx.Vin {
						if in.CanUnlockWith(address) {
							inTxId := hex.EncodeToString(in.TXid)
							spentTXOs[inTxId] = append(spentTXOs[inTxId], in.Vout)
						}
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var utxos []TXOutput
	unspentTranasctions := bc.FindUnspentTransactions(address)
	for _, tx := range unspentTranasctions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				utxos = append(utxos, out)
			}
		}
	}
	return utxos
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}
