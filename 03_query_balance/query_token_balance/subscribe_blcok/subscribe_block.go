package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// 订阅区域块
func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("INFURA_API_KEY")
	client, err := ethclient.Dial("wss://sepolia.infura.io/ws/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex())
			//BlockByHash 需要精确匹配区块哈希，对节点数据同步要求更高
			//BlockByNumber 只需要区块高度，通常更稳定可靠
			//在高速变化的区块链网络中，这种差异是正常的
			block, err := client.BlockByHash(context.Background(), header.Hash()) // client.BlockByHash 貌似查不到数据，改用BlockByNumber
			if block == nil {
				block, err = client.BlockByNumber(context.Background(), header.Number)
			}

			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(header.Nonce.Uint64())
			fmt.Println(header.Time)
			fmt.Println(header.Number.Uint64())

			fmt.Println(block.Hash().Hex())
			fmt.Println(block.Number().Uint64())
			fmt.Println(block.Time())
			fmt.Println(block.Nonce())
			fmt.Println(len(block.Transactions()))
		}
	}
}
