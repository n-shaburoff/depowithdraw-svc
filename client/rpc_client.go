package client

import "encoding/json"

type RpcClient struct {
	endpoint	string
}

func NewRpcClient(enpoint string) *RpcClient {
	return &RpcClient{
		endpoint: enpoint,
	}
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