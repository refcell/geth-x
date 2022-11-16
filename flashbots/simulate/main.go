// Package main demonstrates simulating a transaction using flashbots.
// Uses [metachris/flashbotsrpc]: https://github.com/metachris/flashbotsrpc
package main

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/metachris/flashbotsrpc"
)

func main() {
	// Create a new private key
	privateKey, _ := crypto.GenerateKey()

	// Simulate Transactions using the eth_callBundle rpc method
	rpc := flashbotsrpc.New("https://relay.flashbots.net")
	opts := flashbotsrpc.FlashbotsCallBundleParam{
		Txs:              []string{"YOUR_RAW_TX"},
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
