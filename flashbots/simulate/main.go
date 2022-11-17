// Package main demonstrates simulating a transaction using flashbots.
// Uses [metachris/flashbotsrpc]: https://github.com/metachris/flashbotsrpc
package main

import (
	"errors"
	"fmt"
  "bytes"
  "math/big"
  "strings"

  "github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	// "github.com/ethereum/go-ethereum/ethclient"
	// "github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/metachris/flashbotsrpc"
)

//go:generate abigen --abi ./dai/dai.abi --pkg dai --type Dai --out ./dai/dai.go
type Dai interface {
	Allowance(common.Address, common.Address) (*big.Int, error)
}

func main() {
  // Create RPC Client
	// rpcUrl := os.Getenv("ETH_RPC_URL")
	// rpcClient, err := rpc.Dial(rpcUrl)
	// if err != nil {
	// 	fmt.Printf("Failed to construct rpc client: %v\n", err)
	// 	os.Exit(1)
	// }
	// ethClient := ethclient.NewClient(rpcClient)
	fmt.Println("✅ Configured RPC Client")

	// Create a new private key
	privateKey, _ := crypto.GenerateKey()

	// Create Dai Instance
	// daiContract := 0x6b175474e89094c44da98b954eedeac495271d0f
 //  daiCaller, err := dai.NewDaiTransactor(daiContract, rpcClient)
 //  if err != nil {
 //  	fmt.Printf("Error calling dai: %v\n", err)
 //    os.Exit(1)
 //  }

  // Construct Transaction
  testAddr := common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")
  emptyEip2718Tx := types.NewTx(&types.AccessListTx{
		ChainID:  big.NewInt(1),
		Nonce:    3,
		To:       &testAddr,
		Value:    big.NewInt(10),
		Gas:      25000,
		GasPrice: big.NewInt(1),
		Data:     common.FromHex("5544"),
	})
	signedEip2718Tx, _ := emptyEip2718Tx.WithSignature(
		types.NewEIP2930Signer(big.NewInt(1)),
		common.Hex2Bytes("c9519f4f2b30335884581971573fadf60c6204f59a911df35ee8a540456b266032f1e8e2c5dd761f9e4f88f41c8310aeaba26a8bfcdacfedfa12ec3862d3752101"),
	)   

	// Encode Transaction to bytes
	encodedTx, err := rlp.EncodeToBytes(signedEip2718Tx)
	stringTx := hexutil.Encode(encodedTx)
	if err != nil {
		fmt.Errorf("encode error: %v", err)
	}
	want := common.FromHex("b86601f8630103018261a894b94f5374fce5edbc8e2a8697c15331677e6ebf0b0a825544c001a0c9519f4f2b30335884581971573fadf60c6204f59a911df35ee8a540456b2660a032f1e8e2c5dd761f9e4f88f41c8310aeaba26a8bfcdacfedfa12ec3862d37521")
	if !bytes.Equal(encodedTx, want) {
		fmt.Errorf("encoded RLP mismatch, got %x", encodedTx)
	}

	stringTx = strings.ReplaceAll(stringTx, "0x", "")
	fmt.Printf("Constructed raw tx: %s\n", stringTx)

	// Simulate Transactions using the eth_callBundle rpc method
	rpc := flashbotsrpc.New("https://relay.flashbots.net")
	opts := flashbotsrpc.FlashbotsCallBundleParam{
		Txs:              []string{stringTx},
		BlockNumber:      fmt.Sprintf("0x%x", 13281018),
		StateBlockNumber: "latest",
	}

	result, err := rpc.FlashbotsCallBundle(privateKey, opts)
	if err != nil {
		if errors.Is(err, flashbotsrpc.ErrRelayErrorResponse) { // user/tx error, rather than JSON or network error
			fmt.Println(err.Error())
		} else {
			fmt.Printf("error: %+v\n", err)
		}
		return
	}

	// Print result
	fmt.Printf("%+v\n", result)
}
