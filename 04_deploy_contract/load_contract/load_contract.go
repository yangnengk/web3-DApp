package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/learn/04_deploy_contract/store"
	"log"
	"os"
)

const (
	contractAddress = "0xb68fd10adc559603516fcd149d54045f96e1e87e"
)

// 加载合约
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
	storeContract, err := store.NewStore(common.HexToAddress(""), client)
	if err != nil {
		log.Fatal(err)
	}
	_ = storeContract
}
