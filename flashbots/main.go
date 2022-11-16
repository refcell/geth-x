// Package main demonstrates how to send a bundle to Flashbots.
//
// This follows the example provided in the Flashbots documentation (linked below).
// <https://docs.flashbots.net/flashbots-auction/searchers/quick-start>
//
// When you send bundles to Flashbots they must be signed with a private key so that flashbots can establish identity and reputation for searchers.
// This private key does not store funds and is not the primary private key you use for executing transactions.
// Again, it is only used for identity, and it can be any private key.
//
// Second, you'll need a way to interact with Flashbots.
// The Flashbots builder receives bundles at relay.flashbots.net.
// To send transactions, Flashbots provides specific RPC endpoints.
// For example, below we demonstrate using the glashbots_getUserStats rpc endpoint to fetch the searcher's statistics.
package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	j               = "application/json"
	flashbotURL     = "https://relay.flashbots.net"
	stats           = "flashbots_getUserStats"
	flashbotXHeader = "X-Flashbots-Signature"
	p               = "POST"
)

var (
	privateKey, _ = crypto.HexToECDSA("2e19800fcbbf0abb7cf6d72ee7171f08943bc8e5c3568d1d7420e52136898154")
)

func flashbotHeader(signature []byte, privateKey *ecdsa.PrivateKey) string {
	return crypto.PubkeyToAddress(privateKey.PublicKey).Hex() + ":" + hexutil.Encode(signature)
}

func main() {
	mevHTTPClient := &http.Client{
		Timeout: time.Second * 3,
	}
	currentBlock := big.NewInt(12_900_000)
	params := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  stats,
		"params": []interface{}{
			fmt.Sprintf("0x%x", currentBlock.Uint64()),
		},
	}
	payload, _ := json.Marshal(params)
	req, _ := http.NewRequest(p, flashbotURL, bytes.NewBuffer(payload))
	headerReady, _ := crypto.Sign(
		accounts.TextHash([]byte(hexutil.Encode(crypto.Keccak256(payload)))),
		privateKey,
	)
	req.Header.Add("content-type", j)
	req.Header.Add("Accept", j)
	req.Header.Add(flashbotXHeader, flashbotHeader(headerReady, privateKey))
	resp, _ := mevHTTPClient.Do(req)
	res, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(res))
}
