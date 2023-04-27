package main

import (
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
	config := serverConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}
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
}
