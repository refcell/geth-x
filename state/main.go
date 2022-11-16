package main

import (
	"fmt"
	"math"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

//go:generate abigen --abi ./dai.abi --pkg main --type Dai --out ./dai.go
type Dai interface {
	Allowance(common.Address, common.Address) (*big.Int, error)
}

func main() {
	// Create RPC Client
	rpcUrl := os.Getenv("ETH_RPC_URL")
	rpcClient, err := rpc.Dial(rpcUrl)
	if err != nil {
		fmt.Printf("Failed to construct rpc client: %v\n", err)
		os.Exit(1)
	}
	_ = ethclient.NewClient(rpcClient)
	fmt.Println("✅ Configured RPC Client")

	// Setup Transaction Variables
	fromAddress := common.HexToAddress("0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe")
	toAddress := common.HexToAddress("0x3FB7501f5e451509Da23aD25c331A0737ef514A2")
	allowanceSlot := 3 // Allowance slot (differs from contract to contract)
	dai := common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f")

	// In order to get the allowance slot for the stateDiff, we need to calculate the slot location.
	// The slot location is calculated by first hashing the from address with the allowance slot.
	// Then, we hash the to address with the that result.
	// slot: keccak256(toAddress . keccak256(fromAddress . allowanceSlot))
	// Hash the from address with the slot
	intermediate := crypto.Keccak256(fromAddress.Bytes(), common.BigToHash(big.NewInt(int64(allowanceSlot))).Bytes())
	finalBytes := crypto.Keccak256(toAddress.Bytes(), intermediate)
	finalString := common.Bytes2Hex(finalBytes)
	final := common.HexToHash(finalString)

	fmt.Printf("✅ Calculated slot location: %s\n", final)
	fmt.Printf("Expected: 2b0a4d104c15978ca553d6173c81b852539ee7ea7baee7307410d1b224a172eb\n")

	// Construct dai approval transaction
	// const { data } = await Dai.populateTransaction.allowance(fromAddr, toAddr);

	// eth_call default params
	defaultParams := []interface{}{
		struct {
			From common.Address
			To   common.Address
			Data string
		}{
			From: fromAddress,
			To:   dai,
			Data: "0x",
		},
		"latest",
	}
	fmt.Printf("Using default params: %+v\n", defaultParams)

	// State Diff example
	inner := make(map[common.Hash]common.Hash)
	inner[final] = common.BigToHash(big.NewInt(math.MaxInt64))
	stateDiff := make(map[common.Address]struct{ StateDiff interface{} })
	stateDiff[dai] = struct{ StateDiff interface{} }{inner}

	// Call with no state overrides
	vanillaCall := rpcClient.Call("eth_call", "latest", defaultParams)

	// Call with state overrides
	statefulCall := rpcClient.Call("eth_call", "latest", defaultParams, stateDiff)

	// Print results
	fmt.Printf("✅ Vanilla call: %v\n", vanillaCall)
	fmt.Printf("✅ Stateful call: %v\n", statefulCall)
}
