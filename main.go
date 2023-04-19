package library

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/yunlong-le/library/chaincode"
)

func main() {
	assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating library chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting library chaincode: %v", err)
	}
}
