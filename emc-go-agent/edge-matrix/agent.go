package edge_matrix

import (
	"crypto/ecdsa"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/emc-go-agent/edge-matrix/contracts"
	"github.com/emc-go-agent/edge-matrix/crypto"
	"github.com/emc-go-agent/edge-matrix/helper/rpc"
	"github.com/emc-go-agent/edge-matrix/types"
	"github.com/hashicorp/go-hclog"
	"sync"
)

type EdgeApiMethod string

const METHOD_POST EdgeApiMethod = "POST"
const METHOD_GET EdgeApiMethod = "GET"

type EdgeCallReuslt struct {
	TelegramHash string `json:"telegram_hash"`
	Response     string `json:"response"`
	Err          string `json:"err"`
}

type EdgeAgent struct {
	logger hclog.Logger

	sync.Mutex
	nextNonce        uint64
	nonceCacheEnable bool
	jsonRpcClient    *rpc.JsonRpcClient
	privateKey       *ecdsa.PrivateKey
	address          types.Address
}

func NewDefaultAgent(
	logger hclog.Logger,
	privateKey *ecdsa.PrivateKey,
	jsonRpcClient *rpc.JsonRpcClient) (*EdgeAgent, error) {
	agent := &EdgeAgent{
		logger:        logger,
		jsonRpcClient: jsonRpcClient,
		privateKey:    privateKey,
	}

	address, err := crypto.GetAddressFromKey(agent.privateKey)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to extract key, %v", err.Error()))
	}
	agent.address = address
	return agent, nil
}

func (e *EdgeAgent) GetNoCachedNextNonce() (uint64, error) {
	nonce, err := e.jsonRpcClient.GetNextNonce(e.address.String())
	if err != nil {
		e.logger.Error("unable to get next nonce, %v", err)
		return 0, err
	}
	return nonce, nil
}

func (e *EdgeAgent) GetNextNonce() (uint64, error) {
	e.Lock()
	defer e.Unlock()
	if !e.nonceCacheEnable {
		nonce, err := e.jsonRpcClient.GetNextNonce(e.address.String())
		if err != nil {
			e.logger.Error("unable to get next nonce, %v", err)
			return 0, err
		}
		e.nextNonce = nonce
		e.nonceCacheEnable = true
		return e.nextNonce, nil
	}
	e.nextNonce += 1
	return e.nextNonce, nil
}

func (e *EdgeAgent) DisableNonceCache() {
	e.Lock()
	defer e.Unlock()

	e.nonceCacheEnable = false
}

func (e *EdgeAgent) CallEdgeInfo(nodeId string) (*EdgeCallReuslt, error) {
	nonce, err := e.GetNextNonce()
	if err != nil {
		return nil, err
	}
	input := `{"peerId":"%s","endpoint":"/info","Input":{}}`
	response, err := e.jsonRpcClient.SendRawTelegram(
		contracts.EdgeCallPrecompile,
		nonce,
		fmt.Sprintf(input, nodeId),
		e.privateKey,
	)
	if err != nil {
		return nil, err
	}
	result := EdgeCallReuslt{}
	result.TelegramHash = response.Result.TelegramHash
	respBytes, err := base64.StdEncoding.DecodeString(response.Result.Response)
	if err != nil {
		e.logger.Debug("CallEdgeApi -->base64 decode", "err", err.Error())
		result.Err = fmt.Sprintf("base64 decode err: %s", err.Error())
		return &result, nil
	}
	result.Response = string(respBytes)
	return &result, nil
}

func (e *EdgeAgent) CallEdgeApi(nodeId, path, data string, method EdgeApiMethod) (*EdgeCallReuslt, error) {
	nonce, err := e.GetNextNonce()
	if err != nil {
		return nil, err
	}
	input := `{"peerId":"%s","endpoint":"/api","Input":{"method": "%s","headers":[],"path":"%s","body":%s}}`
	response, err := e.jsonRpcClient.SendRawTelegram(
		contracts.EdgeCallPrecompile,
		nonce,
		fmt.Sprintf(input, nodeId, method, path, data),
		e.privateKey,
	)
	if err != nil {
		return nil, err
	}
	result := EdgeCallReuslt{}
	result.TelegramHash = response.Result.TelegramHash
	respBytes, err := base64.StdEncoding.DecodeString(response.Result.Response)
	if err != nil {
		e.logger.Debug("CallEdgeApi -->base64 decode", "err", err.Error())
		result.Err = fmt.Sprintf("base64 decode err: %s", err.Error())
		return &result, nil
	}
	result.Response = string(respBytes)
	return &result, nil
}
