package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"log"
)

// 创建一个钱包
func main() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	//privateKey, err := crypto.HexToECDSA("ccec5314acec3d18eae81b6bd988b844fc4f7f7d3c828b351de6d0fede02d3f2")
	//if err != nil {
	//	log.Fatal(err)
	//}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes[2:])) //去掉前两个字节 ‘0x’
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("from pubkey:", hexutil.Encode(publicKeyBytes[4:])) //去掉前四个字节 ‘0x04’
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("address short:", address)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println("full: ", hexutil.Encode(hash.Sum(nil)[:]))    // nil 是最常用的模式，表示"只需要哈希结果",避免了创建空切片的步骤
	fmt.Println("short: ", hexutil.Encode(hash.Sum(nil)[12:])) // crypto.PubkeyToAddress 和 此处结果一样
}
