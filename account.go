package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
)

// Account 描述账户
type Account struct {
	ID        string
	DocType   string
	MSPID     string
	PublicKey string
	Algorithm string
}

func paramsAccountFromJSON(v []byte) (*Account, error) {
	var account Account
	if err := json.Unmarshal(v, &account); err != nil {
		return nil, errors.Wrap(err, "解析Account出错")
	}
	return &account, nil
}

// JSON 返回JSON编码的字节数组
func (a Account) JSON() []byte {
	buf, _ := json.Marshal(a)
	return buf
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	PublicKey string
	Algorithm string
}

// JSON 转换为JSON格式
func (r RegisterRequest) JSON() []byte {
	buf, _ := json.Marshal(r)
	return buf
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	ID string
}

// JSON 返回JSON编码的字节数组
func (r RegisterResponse) JSON() []byte {
	buf, _ := json.Marshal(r)
	return buf
}

// Register 注册账户
func Register(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if err := checkParamsLength(params, 1); err != nil {
		return shim.Error(err.Error())
	}
	var req RegisterRequest
	if err := json.Unmarshal([]byte(params[0]), &req); err != nil {
		return shim.Error(err.Error())
	}
	clientID, err := cid.New(stub)
	if err != nil {
		err = errors.Wrap(err, "获取CID失败")
		return shim.Error(err.Error())
	}

	id, err := generateID(clientID)
	if err != nil {
		return shim.Error(err.Error())
	}

	mspid, err := clientID.GetMSPID()
	if err != nil {
		err = errors.Wrap(err, "获取MSPID失败")
		return shim.Error(err.Error())
	}

	var account = Account{
		ID:        id,
		MSPID:     mspid,
		DocType:   "Account",
		PublicKey: req.PublicKey,
		Algorithm: req.Algorithm,
	}

	buf, err := json.Marshal(account)
	if err != nil {
		err = errors.Wrap(err, "编码记录失败")
		return shim.Error(err.Error())
	}
	ck, err := stub.CreateCompositeKey("Account", []string{account.ID})
	if err != nil {
		err = errors.Wrap(err, "创建复合主键失败")
		return shim.Error(err.Error())
	}

	if err := stub.PutState(ck, buf); err != nil {
		err = errors.Wrap(err, "保存记录失败")
		return shim.Error(err.Error())
	}

	resp := &RegisterResponse{ID: id}
	return shim.Success(resp.JSON())
}

// GetAccount 获取用户信息
func GetAccount(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if err := checkParamsLength(params, 1); err != nil {
		return shim.Error(err.Error())
	}

	id := params[0]
	if id == "" {
		err := errors.New("id不能为空")
		return shim.Error(err.Error())
	}

	if strings.ToLower(id) == "self" {
		clientID, err := cid.New(stub)
		if err != nil {
			err = errors.Wrap(err, "获取CID失败")
			return shim.Error(err.Error())
		}
		id, err = generateID(clientID)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	ck, err := stub.CreateCompositeKey("Account", []string{id})
	if err != nil {
		err = errors.Wrap(err, "创建复合键出错")
		return shim.Error(err.Error())
	}

	state, err := stub.GetState(ck)
	if err != nil {
		err = errors.Wrap(err, "获取记录失败")
		return shim.Error(err.Error())
	}
	if state == nil {
		err = errors.Errorf("找不到ID: %s 对应的记录", id)
		return shim.Error(err.Error())
	}
	return shim.Success(state)
}

func generateID(clientID *cid.ClientID) (string, error) {
	crtID, err := clientID.GetID()
	if err != nil {
		return "", errors.Wrap(err, "获取ID失败")
	}
	hash := sha256.New()
	if _, err := hash.Write([]byte(crtID)); err != nil {
		return "", errors.Wrap(err, "创建id失败")
	}
	return base64.RawStdEncoding.EncodeToString(hash.Sum(nil)), nil
}
