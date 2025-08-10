package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
	"time"
)

const (
	contractAddress = "0xb68fd10adc559603516fcd149d54045f96e1e87e"
)

// 不使用abi执行合约
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
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// ------------------------------------------------------------------
	//gasPrice, err := client.SuggestGasPrice(context.Background())
	//gasPrice = big.NewInt(0).Add(gasPrice, big.NewInt(1000000))
	//if err != nil {
	//	log.Fatal(err)
	//}
	gasTipCap, err := client.SuggestGasTipCap(context.Background()) //用户愿意支付给矿工/验证者的最大小费（优先费用）
	if err != nil {
		log.Fatal(err)
	}
	head, err := client.HeaderByNumber(context.Background(), nil) //用户愿意为交易支付的总费用上限
	if err != nil {
		log.Fatal(err)
	}
	gasFeeCap := new(big.Int).Add(head.BaseFee, new(big.Int).Mul(gasTipCap, big.NewInt(2)))
	// --------------------------------------------------------------------------------------------
	// 准备交易数据
	//contractAbi, err := abi.JSON(strings.NewReader("[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_version\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"ItemSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"items\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"setItem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//methodName := "setItem"
	//var key [32]byte
	//var value [32]byte
	//copy(key[:], []byte("demo_save_key_use_abi"))
	//copy(value[:], []byte("demo_save_value_use_abi_1111"))
	//input, err := contractAbi.Pack(methodName, key, value)
	methodSignature := []byte("setItem(bytes32,bytes32)")
	methodSelector := crypto.Keccak256(methodSignature)[:4]
	var key [32]byte
	var value [32]byte
	copy(key[:], []byte("demo_save_key_no_use_abi"))
	copy(value[:], []byte("demo_save_value_no_use_abi_1"))
	// 组合调用数据
	var input []byte
	input = append(input, methodSelector...)
	input = append(input, key[:]...)
	input = append(input, value[:]...)
	// --------------------------------------------------------------------
	// 创建交易并签名
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	contractAdd := common.HexToAddress(contractAddress)
	// -----------------------------------------------------
	//types.NewTransaction()
	/*
		LegacyTx (types.LegacyTx):
			使用传统的交易格式
			包含 GasPrice 字段，指定愿意为每单位 gas 支付的固定价格
			结构简单，兼容所有 Ethereum 网络版本
		DynamicFeeTx (types.DynamicFeeTx):
			为 EIP-1559 升级引入的新交易类型
			使用 GasTipCap (maxPriorityFeePerGas) 和 GasFeeCap (maxFeePerGas) 两个字段
			GasTipCap: 给矿工的小费上限
			GasFeeCap: 总费用上限 (基础费用 + 小费)
	*/
	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		To:        &contractAdd,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       310000,
		Value:     big.NewInt(0),
		Data:      input,
	})
	//tx := types.NewTx(&types.LegacyTx{
	//	Nonce:    nonce,
	//	To:       &contractAdd,
	//	GasPrice: gasPrice,
	//	Gas:      310000,
	//	Value:    big.NewInt(0),
	//	Data:     input,
	//})
	// ---------------------------------------------------------------------------------
	/*
		Ethereum 中用于交易签名的两种不同签名器（signer），它们的主要区别如下：
		EIP155Signer (types.NewEIP155Signer(chainID))：
		实现了 EIP-155 规范，该规范引入了重放攻击保护机制
		在签名交易时会将 chainID 包含在签名过程中，防止交易在不同网络间被重放
		适用于较老的交易类型（如 Legacy transactions）
		签名格式：r, s, v，其中 v 包含了 chainID 信息
		LondonSigner (types.NewLondonSigner(chainID))：
		为伦敦升级（London upgrade）后的网络设计，支持新的交易类型（如 EIP-1559 交易）
		兼容 EIP-155，但可以处理更现代的交易格式，包括动态费用交易（Dynamic Fee Transactions）
		支持所有交易类型：Legacy、AccessList 和 DynamicFee transactions
		对于 EIP-1559 交易，能正确处理 gasTipCap 和 gasFeeCap 字段
		在你的代码中，由于使用的是 LegacyTx 类型的交易，两种签名器都能正常工作。但如果使用 DynamicFeeTx（EIP-1559 交易），则必须使用 LondonSigner 或更新的签名器。
	*/
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("transaction send %s\n", signedTx.Hash().Hex())
	_, err = waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// 查询刚刚设置的值
	//callInput, err := contractAbi.Pack("items", key)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//callMsg := ethereum.CallMsg{
	//	To:   &contractAdd,
	//	Data: callInput,
	//}
	itemSignature := []byte("items(bytes32)")
	itemsSelector := crypto.Keccak256(itemSignature)[:4]
	var callInput []byte
	callInput = append(callInput, itemsSelector...)
	callInput = append(callInput, key[:]...)
	callMsg := ethereum.CallMsg{
		To:   &contractAdd,
		Data: callInput,
	}

	// 解析返回值
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}
	var unpackedValue [32]byte
	//copy(unpackedValue[:], result)	// 和下面的代码效果一样
	//err = contractAbi.UnpackIntoInterface(&unpackedValue, "items", result)
	//if err != nil {
	//	log.Fatal(err)
	//}
	copy(unpackedValue[:], result)
	fmt.Println("is value in contract equals to origin value:", unpackedValue == value)
}

func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if receipt != nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		time.Sleep(time.Second * 1)
	}
}
