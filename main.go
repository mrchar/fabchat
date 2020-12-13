package main

import "github.com/hyperledger/fabric-chaincode-go/shim"

func main() {
	if err := shim.Start(new(Chaincode)); err != nil {
		panic(err)
	}
}
