/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"CRUD-Hyperledger-chaincode/chaincode"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	votingSmartContract, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating auction chaincode: %v", err)
	}

	if err := votingSmartContract.Start(); err != nil {
		log.Panicf("Error starting auction chaincode: %v", err)
	}
}
