package Blockchain_Go

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

type Transaction struct {
	Amount   float64
	Sender   string
	Reciever string
}

func (transaction *Transaction) sign() string {
	signature := sha256.New()
	signature.Write([]byte(transaction.Sender))
	signature.Write([]byte(transaction.Reciever))
	signature.Write([]byte(fmt.Sprint(transaction.Amount)))

	signatureHash := string(signature.Sum(nil))

	return signatureHash
}

func (transaction *Transaction) verify(signature string) bool {
	if signature == transaction.sign() {
		return true
	} else {
		return false
	}
}

type Block struct {
	Index        int
	TimeStamp    string
	T            Transaction
	PreviousHash string
	Nonce        int
	Hash         string
}

func (block *Block) computeHash() string {
	hash := sha256.New()

	hash.Write([]byte(fmt.Sprint(block.Index)))
	hash.Write([]byte((block.TimeStamp)))
	hash.Write([]byte(fmt.Sprint(block.T)))
	hash.Write([]byte(fmt.Sprint(block.Nonce)))
	hash.Write([]byte(block.PreviousHash))

	return string(hash.Sum(nil))
}

func (block *Block) proofOfWork(difficulty int) {

	for block.Hash[0:difficulty] != strings.Repeat("0", difficulty+1) {
		block.Nonce++
		block.Hash = block.computeHash()
	}

}

type Blockchain struct {
	currency   string
	chain      []Block
	difficulty int
}

func (blockchain *Blockchain) createBlockchain() {
	blockchain.chain = append(blockchain.chain, Block{0, time.Now().String(), Transaction{0, "1", "2"}, "0", 0, "0"})
	blockchain.chain[0].Hash = blockchain.chain[0].computeHash()
}

func (blockchain *Blockchain) getLatestBlock() Block {
	return blockchain.chain[len(blockchain.chain)-1]
}

func (blockchain *Blockchain) addBlock(block Block) {
	block.PreviousHash = blockchain.getLatestBlock().Hash
	block.proofOfWork(blockchain.difficulty)
	blockchain.chain = append(blockchain.chain, block)
}

func (blockchain *Blockchain) isChainValid() bool {
	for i := 1; i < len(blockchain.chain); i++ {
		currentBlock := blockchain.chain[i]
		previousBlock := blockchain.chain[i-1]

		if currentBlock.Hash != currentBlock.computeHash() {
			return false
		}
		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}

	return true
}

func (blockchain *Blockchain) pushTransactionToBlock(transaction Transaction) bool {
	signature := transaction.sign()
	if transaction.verify(signature) {
		blockchain.addBlock(Block{blockchain.getLatestBlock().Index + 1, time.Now().String(), transaction, *new(string), 0, *new(string)})
		return true
	} else {
		return false
	}
}

func (blockchain *Blockchain) getWalletBalance(address string) float64 {
	balance := 0.0

	if address == blockchain.getLatestBlock().T.Sender {
		balance -= blockchain.getLatestBlock().T.Amount
	}
	if address == blockchain.getLatestBlock().T.Reciever {
		balance += blockchain.getLatestBlock().T.Amount
	}

	return balance
}

type Wallet struct {
	TotalBalance float64
	Address      string
	Currency     Blockchain
}

func (wallet *Wallet) updateBalance() {
	wallet.TotalBalance += wallet.Currency.getWalletBalance(wallet.Address)
}
