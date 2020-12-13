package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
)

// Message 消息
type Message struct {
	ID        string
	From      string
	To        string
	Timestamp int64
	Encrypted bool
	PubKey    string `json:",omitempty"`
	Algorithm string `json:",omitempty"`
	Content   string
}

// JSON 返回JSON编码的字节数组
func (m Message) JSON() []byte {
	buf, _ := json.Marshal(m)
	return buf
}

// SendRequest 发送消息的请求
type SendRequest struct {
	To        string
	Encrypted bool
	PubKey    string `json:",omitempty"`
	Algorithm string `json:",omitempty"`
	Content   string
}

func parseSendRequest(v []byte) (*SendRequest, error) {
	var req SendRequest
	if err := json.Unmarshal(v, &req); err != nil {
		return nil, errors.Wrap(err, "解析接受到的请求出错")
	}
	return &req, nil
}

// JSON 返回JSON编码的字节数组
func (r SendRequest) JSON() []byte {
	buf, _ := json.Marshal(r)
	return buf
}

// SendResponse 发送请求的相应
type SendResponse struct {
	ID string
}

// JSON 返回JSON字节数组
func (r SendResponse) JSON() []byte {
	buf, _ := json.Marshal(r)
	return buf
}

// Send 发送消息
func Send(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	// 检查参数
	if err := checkParamsLength(params, 1); err != nil {
		return shim.Error(err.Error())
	}

	// 解析请求
	req, err := parseSendRequest([]byte(params[0]))
	if err != nil {
		return shim.Error(err.Error())
	}

	// 检查请求
	if req.Encrypted && (req.PubKey == "" || req.Algorithm == "") {
		return shim.Error("加密消息必须指明加密公钥和加密算法")
	}

	// 获取Sender
	clientID, err := cid.New(stub)
	if err != nil {
		err = errors.Wrap(err, "创建clientID失败")
		return shim.Error(err.Error())
	}

	sender, err := generateID(clientID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// 创建ID
	hash := sha256.New()
	if _, err := hash.Write([]byte(stub.GetTxID() + "::" + sender + "::" + req.To)); err != nil {
		return shim.Error("创建ID失败")
	}
	id := base64.RawStdEncoding.EncodeToString(hash.Sum(nil))

	// 获取时间戳
	stamp, err := stub.GetTxTimestamp()
	if err != nil {
		err = errors.Wrap(err, "获取时间戳出错")
		return shim.Error(err.Error())
	}

	// 构建消息
	var message = Message{
		ID:        id,
		From:      sender,
		To:        req.To,
		Timestamp: stamp.Seconds,
		Encrypted: req.Encrypted,
		PubKey:    req.PubKey,
		Algorithm: req.Algorithm,
		Content:   req.Content,
	}

	// 创建复合键
	ck, err := stub.CreateCompositeKey("Message", []string{id})
	if err != nil {
		err = errors.Wrap(err, "创建复合主键出错")
		return shim.Error(err.Error())
	}

	// 保存记录
	if err := stub.PutState(ck, message.JSON()); err != nil {
		err = errors.Wrap(err, "保存记录出错")
		return shim.Error(err.Error())
	}

	// 设置事件
	if err := stub.SetEvent("New Message To "+message.To, []byte(message.ID)); err != nil {
		err = errors.Wrap(err, "设置事件失败")
		return shim.Error(err.Error())
	}

	return shim.Success(SendResponse{ID: message.ID}.JSON())
}

// Receive 接受消息
func Receive(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	// 检查参数
	if err := checkParamsLength(params, 1); err != nil {
		return shim.Error(err.Error())
	}

	// 创建复合键
	ck, err := stub.CreateCompositeKey("Message", []string{params[0]})
	if err != nil {
		err = errors.Wrap(err, "创建复合主键出错")
		return shim.Error(err.Error())
	}

	state, err := stub.GetState(ck)
	if err != nil {
		err = errors.Wrap(err, "读取消息失败")
	}

	if state == nil {
		err = errors.Errorf("找不到ID: %s 对应的记录", params[0])
		return shim.Error(err.Error())
	}

	return shim.Success(state)
}
