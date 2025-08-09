package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	token "github.com/learn/03_query_balance/query_token_balance/erc20"
	"log"
	"math"
	"math/big"
	"os"
)

// 查询代币余额
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
	// Golem (GNT) Address 代币合约地址
	tokenAddress := common.HexToAddress("0xE165AE29c619F455acD5F0b1C81958C229fCc0B0")
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	// 账户地址
	accountAddress := common.HexToAddress("0x3fA4B0E71d2e042C82d4532d8D320D1caa765026") // 0x25836239F76632635F815689389C537133248edb
	bal, err := instance.BalanceOf(&bind.CallOpts{}, accountAddress)
	if err != nil {
		log.Fatal(err)
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	decimal, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Name:", name, "Symbol:", symbol, "Decimal:", decimal, "wei Balance:", bal)
	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimal))))
	fmt.Println("eth balance:", value)
}
