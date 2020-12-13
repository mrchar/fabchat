package main

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

func TestSendMessage(t *testing.T) {
	stub := newMockStub(SerializedIdentity1)

	// 注册
	var regReq = RegisterRequest{
		PublicKey: "pub",
		Algorithm: "alg",
	}
	regResp, err := mockRegister(stub, regReq)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("register response: %+v", regResp)

	// 发送消息
	var sendReq = SendRequest{
		To:        regResp.ID,
		Encrypted: false,
		Content:   "hello",
	}
	sendResp, err := mockSend(stub, sendReq)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("send response:%+v", sendResp)

	// 接受消息
	message, err := mockReceive(stub, sendResp.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("received message: %+v", message)
}

func mockSend(stub *shimtest.MockStub, req SendRequest) (*SendResponse, error) {
	resp := stub.MockInvoke(
		uuid.New().String(),
		[][]byte{[]byte("Send"), req.JSON()},
	)
	payload, err := parseResp(resp)
	if err != nil {
		return nil, err
	}
	var sendResp SendResponse
	if err := json.Unmarshal(payload, &sendResp); err != nil {
		return nil, err
	}
	return &sendResp, nil
}

func mockReceive(stub *shimtest.MockStub, id string) (*Message, error) {
	resp := stub.MockInvoke(
		uuid.New().String(),
		[][]byte{[]byte("Receive"), []byte(id)},
	)
	payload, err := parseResp(resp)
	if err != nil {
		return nil, err
	}
	var message Message
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}
	return &message, nil
}
