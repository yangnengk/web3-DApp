package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// 查询区块头
func main() {
	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("INFURA_API_KEY")
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}

	blockNumber := big.NewInt(5671744)
	head, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(head.Number.Uint64())
	fmt.Println(head.Time)
	fmt.Println(head.Difficulty.Uint64()) //0
	fmt.Println(head.Hash().Hex())

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	fmt.Println(block.Number().Uint64())
	fmt.Println(block.Time())
	fmt.Println(block.Difficulty().Uint64())
	fmt.Println(block.Hash().Hex())
	fmt.Println(len(block.Transactions()))

	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)

}
