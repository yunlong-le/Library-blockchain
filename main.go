package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/yunlong-le/library/chaincode-2"
)

type serverConfig struct {
	CCID    string
	Address string
}

func main() {
	fmt.Printf("------------------------------------------------------------main start------------------------------------------------------------: %v\n")
	mode := os.Getenv("MODE")
	fmt.Printf("----------------------------------------get MODE: %v\n", mode)
	if mode == "chaincode" {
		err := shim.Start(new(chaincode.SmartContract))
		if err != nil {
			fmt.Printf("-----------------------------------Error starting Simple chaincode: %s", err)
		}
	} else if mode == "external" {
		// External mode startup logic
		config := serverConfig{
			CCID:    os.Getenv("CHAINCODE_ID"),
			Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
		}
		//chaincode, err := contractapi.NewChaincode(&SmartContract{})
		cc := new(chaincode.SmartContract)
		server := &shim.ChaincodeServer{
			CCID:    config.CCID,
			Address: config.Address,
			CC:      cc,
			TLSProps: shim.TLSProperties{
				Disabled: true,
			},
		}
		if err := server.Start(); err != nil {
			log.Panicf("Error starting asset-transfer-basic chaincode: %s", err.Error())
		}
	} else {
		log.Panic("---------------------------------------Invalid mode specified")
	}
}
