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
	mode := os.Getenv("MODE")
	if mode == "chaincode" {
		// 链码模式启动逻辑
		err := shim.Start(new(chaincode.SmartContract))
		if err != nil {
			fmt.Printf("Error starting Simple chaincode: %s", err)
		}
	} else if mode == "external" {
		// 外部模式启动逻辑
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
		log.Panic("Invalid mode specified")
	}
}
