package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
	"strings"
)

var key = "0x0000000000000000000000000000000000000000000000000000000000002222"
var value = "0x0000000000000000000000000000000000000000000000000000000000002000"

var storeAbi = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

var contractAddress = "0xD6620bf645441c593E06F57d0b37d5148eeB31b0"

// 合约事件
func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("INFURA_API_KEY")
	//client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + apiKey)
	client, err := ethclient.Dial("wss://sepolia.infura.io/ws/v3/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}

	//header, err := client.HeaderByNumber(context.Background(), nil)
	//fromBlock := new(big.Int).Sub(header.Number, big.NewInt(10000))

	contractAdd := common.HexToAddress(contractAddress)
	// 查询合约事件
	//queryEvent(contractAdd, client)
	// 订阅合约事件
	subscribeEvent(contractAdd, client)
}

// 订阅合约事件
func subscribeEvent(contractAdd common.Address, client *ethclient.Client) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			contractAdd,
		},
	}
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	contractAbi, err := abi.JSON(strings.NewReader(storeAbi))
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog.BlockHash.Hex())
			fmt.Println(vLog.BlockNumber)
			fmt.Println(vLog.TxHash.Hex())
			event := struct {
				Key   [32]byte
				Value [32]byte
			}{}
			err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(common.Bytes2Hex(event.Key[:]))
			fmt.Println(common.Bytes2Hex(event.Value[:]))
			var topics []string
			for i := range vLog.Topics {
				topics = append(topics, vLog.Topics[i].Hex())
			}
			fmt.Println("topics[0] =", topics[0])
			if len(topics) > 1 {
				fmt.Println("indexed topics:", topics[1:])
			}
		}
	}
}

// 查询合约事件
func queryEvent(contractAdd common.Address, client *ethclient.Client) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6920583),
		//FromBlock: fromBlock,
		// ToBlock:   big.NewInt(2394201),
		Addresses: []common.Address{
			contractAdd,
		},
		// Topics: [][]common.Hash{
		//  {},
		//  {},
		// },
	}
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	contractAbi, err := abi.JSON(strings.NewReader(storeAbi))
	if err != nil {
		log.Fatal(err)
	}
	for _, vLog := range logs {
		fmt.Println(vLog.BlockHash.Hex())
		fmt.Println(vLog.BlockNumber)
		fmt.Println(vLog.TxHash.Hex())
		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(common.Bytes2Hex(event.Key[:]))
		fmt.Println(common.Bytes2Hex(event.Value[:]))
		var topics []string
		for i := range vLog.Topics {
			topics = append(topics, vLog.Topics[i].Hex())
		}
		fmt.Println("topics[0] =", topics[0])
		if len(topics) > 1 {
			fmt.Println("indexed topics:", topics[1:])
		}
	}
	eventSignature := []byte("ItemSet(bytes32,bytes32)")
	hash := crypto.Keccak256Hash(eventSignature)
	fmt.Println("signature topics =", hash.Hex())
}
