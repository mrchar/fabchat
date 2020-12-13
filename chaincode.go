package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// Chaincode 实现链码
type Chaincode struct{}

// Init 执行初始化请求
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke 执行调用请求
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fcn, params := stub.GetFunctionAndParameters()
	switch fcn {
	case "Register":
		return Register(stub, params)
	case "GetAccount":
		return GetAccount(stub, params)
	case "Send":
		return Send(stub, params)
	case "Receive":
		return Receive(stub, params)
	default:
		return shim.Error(fmt.Sprintf("未知的方法: %s", fcn))
	}
}
