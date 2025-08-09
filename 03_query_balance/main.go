package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"log"
	"math"
	"math/big"
	"os"
)

// 查询账户余额
func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("INFURA_API_KEY")
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	latestBlockNumber, err := client.BlockNumber(context.Background())
	println(latestBlockNumber)
	account := common.HexToAddress("0x3fA4B0E71d2e042C82d4532d8D320D1caa765026")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(balance)
	blockNumber := big.NewInt(8945500)
	balance2, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance2)
	fBalance := new(big.Float)
	fBalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18))) // 执行除法 除以10的18次方
	fmt.Println(ethValue)

	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println(pendingBalance)
}
