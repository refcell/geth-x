package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	// Create RPC Client
	rpcUrl := os.Getenv("ETH_RPC_URL")
	rpcClient, err := rpc.Dial(rpcUrl)
	if err != nil {
		fmt.Printf("Failed to construct rpc client: %v\n", err)
		os.Exit(1)
	}
	ethClient := ethclient.NewClient(rpcClient)
	fmt.Println("✅ Configured RPC Client")

	// Get a Transaction Receipt
	txHashString := "0x818b700de19d807ff9fa23c52fd6afd9dad7a419e0b9e6124edeecb4f53cbca8"
	txHash := common.HexToHash(txHashString)
	tx, pending, err := ethClient.TransactionByHash(context.Background(), txHash)
	if err != nil {
		fmt.Printf("Failed to get tx with hash: %s\n", txHashString)
		os.Exit(1)
	}
	fmt.Printf("✅ Fetched transaction for hash: %s\n", txHashString)
	fmt.Printf("Is tx pending: %t\n", pending)
	fmt.Printf("Got tx by hash: %+v\n", tx)

	// Execute the debug trace call
	// Format: debug.TraceCall(args ethapi.CallArgs, blockNrOrHash rpc.BlockNumberOrHash, config *TraceConfig) (*ExecutionResult, error)

	rpcClient.Call(
		"debug_traceCall",
		"14586706",
	// {
	//     from: txResp.from,
	//     to: txResp.to,
	//     value: toRpcHexString(txResp.value),
	//     gas: toRpcHexString(txResp.gasLimit),
	//     data: txResp.data,
	//   },
	//   "14586706",
	//   {
	//     tracer: `{
	//         data: [],
	//         fault: function(log) {},
	//         step: function(log) {
	//             var s = log.op.toString();
	//             if(s == "LOG0" || s == "LOG1" || s == "LOG2" || s == "LOG3" || s == "LOG4") {
	//                 var myStack = [];
	//                 var stackLength = log.stack.length();
	//                 for (var i = 0; i < stackLength; i++) {
	//                     myStack.unshift(log.stack.peek(i));
	//                 }
	//
	//                 var offset = parseInt(myStack[stackLength - 1]);
	//                 var length = parseInt(myStack[stackLength - 2]);
	//                 this.data.push({
	//                     op: s,
	//                     address: log.contract.getAddress(),
	//                     caller: log.contract.getCaller(),
	//                     stack: myStack,
	//                     memory: log.memory.slice(offset, offset + length),
	//                 });
	//             }
	//         },
	//         result: function() { return this.data; }}
	//     `,
	//   },
	)
}
