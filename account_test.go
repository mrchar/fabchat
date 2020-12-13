package main

import (
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
)

var (
	SerializedIdentity1 = &msp.SerializedIdentity{
		Mspid: "org1",
		IdBytes: []byte(`-----BEGIN CERTIFICATE-----
MIICGDCCAb+gAwIBAgIQF/p1+SdXrBrjU66b+NLM2DAKBggqhkjOPQQDAjBzMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eu
b3JnMS5leGFtcGxlLmNvbTAeFw0yMDEyMDYxMjQzMDBaFw0zMDEyMDQxMjQzMDBa
MFsxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1T
YW4gRnJhbmNpc2NvMR8wHQYDVQQDDBZBZG1pbkBvcmcxLmV4YW1wbGUuY29tMFkw
EwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAESKVcwhxVTQbCfBEcBzBp6VzCVOph/mCL
hxB9u069+e1Shhcq6bfLASq8d1Hur5BIePjJk3jDc4ZkBnb2S/jqMKNNMEswDgYD
VR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYDVR0jBCQwIoAgnsNWynDGP1Vx
vqQe0Td/5FtkB7AZGGmwgrW8xWOyzmkwCgYIKoZIzj0EAwIDRwAwRAIgZj8n/hRC
6bRxW8iZkoZP3UmymCNMFukcyIUMQBLQILwCIFqxv132QXbbWSVh+3s6d7X/t0rT
d4dVqqHIg10gI8/i
-----END CERTIFICATE-----
`),
	}
)

func TestRegister(t *testing.T) {
	stub := newMockStub(SerializedIdentity1)

	// 注册
	var req = RegisterRequest{
		PublicKey: "pub",
		Algorithm: "alg",
	}
	regResp, err := mockRegister(stub, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resp: %+v", regResp)

	// 获取
	account, err := mockGetAccount(stub, "self")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("resp: %+v", account)
}

func mockRegister(stub *shimtest.MockStub, req RegisterRequest) (*RegisterResponse, error) {
	resp := stub.MockInvoke(
		uuid.New().String(),
		[][]byte{[]byte("Register"), req.JSON()},
	)

	payload, err := parseResp(resp)
	if err != nil {
		return nil, err
	}

	var regResp RegisterResponse
	if err := json.Unmarshal(payload, &regResp); err != nil {
		return nil, err
	}
	return &regResp, nil
}

func mockGetAccount(stub *shimtest.MockStub, id string) (*Account, error) {
	resp := stub.MockInvoke(
		uuid.New().String(),
		[][]byte{[]byte("GetAccount"), []byte("self")},
	)

	payload, err := parseResp(resp)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := json.Unmarshal(payload, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func newMockStub(sid *msp.SerializedIdentity) *shimtest.MockStub {
	stub := shimtest.NewMockStub("cc", new(Chaincode))
	stub.Creator, _ = proto.Marshal(sid)
	return stub
}

func parseResp(resp peer.Response) ([]byte, error) {
	if resp.Status != shim.OK {
		return nil, errors.New(resp.Message)
	}
	return resp.Payload, nil
}
