package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type RpcClient struct {
	endpoint	string
}

func NewRpcClient(endpoint string) *RpcClient {
	return &RpcClient{
		endpoint: endpoint,
	}
}

func (c *RpcClient) rpcCall(method string, params interface{}) (RpcResponse, error) {
	body, err := json.Marshal(RpcRequest{
		Version: "2.0",
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return RpcResponse{}, errors.Wrap(err, "failed to marshal json")
	}

	resp, err := http.Post(c.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return RpcResponse{}, fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RpcResponse{}, fmt.Errorf("failed to get response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return RpcResponse{}, fmt.Errorf("request failed, status code - %d, response - %s", resp.StatusCode, string(b))
	}

	var rpcResponse RpcResponse
	err = json.Unmarshal(b, &rpcResponse)
	if err != nil {
		return RpcResponse{}, fmt.Errorf("failed to parse response body: %w", err)
	}

	if rpcResponse.Error != nil {
		return rpcResponse, fmt.Errorf("rpc call failed, code - %d, message - %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}

	return rpcResponse, nil
}

type RpcRequest struct {
	Version	string		`json:"jsonrpc"`
	Id 		string		`json:"id"`
	Method	string		`json:"method"`
	Params	interface{}	`json:"params"`
}

type RpcResponse struct {
	Version	string			`json:"jsonrpc"`
	Id 		string			`json:"id"`
	Result 	json.RawMessage	`json:"result"`
	Error	*RpcError		`json:"error"`
}

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AddressResponse struct {
	Address string	`json:"address"`
}

type AddressResult struct {
	Address AddressResponse	`json:"address"`
}