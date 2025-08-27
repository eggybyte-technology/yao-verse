package server

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/eggybyte-technology/yao-portal/interfaces"
	sharedinterfaces "github.com/eggybyte-technology/yao-verse-shared/interfaces"
	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// JSONRPCServiceImpl implements the JSON-RPC service interface
type JSONRPCServiceImpl struct {
	mu               sync.RWMutex
	config           *interfaces.PortalConfig
	txPoolService    interfaces.TxPoolService
	stateReader      sharedinterfaces.StateReader
	messageProducer  sharedinterfaces.MessageProducer
	requestValidator interfaces.RequestValidator
	metrics          *interfaces.JSONRPCMetrics
	startTime        time.Time
}

// NewJSONRPCService creates a new JSON-RPC service instance
func NewJSONRPCService(config *interfaces.PortalConfig) (*JSONRPCServiceImpl, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	service := &JSONRPCServiceImpl{
		config:    config,
		startTime: time.Now(),
		metrics: &interfaces.JSONRPCMetrics{
			RequestsByMethod: make(map[string]uint64),
			ErrorsByType:     make(map[string]uint64),
		},
	}

	return service, nil
}

// Transaction related methods

// SendRawTransaction sends a raw transaction to the network
func (j *JSONRPCServiceImpl) SendRawTransaction(ctx context.Context, data hexutil.Bytes) (common.Hash, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	// Update metrics
	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_sendRawTransaction"]++

	// TODO: Implement actual transaction processing
	// 1. Parse transaction data
	// 2. Validate transaction
	// 3. Add to transaction pool
	// 4. Send to RocketMQ for sequencing

	// For now, return a dummy hash
	return common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"), nil
}

// GetTransactionByHash returns a transaction by its hash
func (j *JSONRPCServiceImpl) GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_getTransactionByHash"]++

	// TODO: Implement transaction lookup
	// 1. Check transaction pool first
	// 2. Query from state reader
	return nil, fmt.Errorf("not implemented")
}

// GetTransactionReceipt returns a transaction receipt by hash
func (j *JSONRPCServiceImpl) GetTransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_getTransactionReceipt"]++

	// TODO: Implement receipt lookup
	return nil, fmt.Errorf("not implemented")
}

// GetTransactionCount returns the transaction count (nonce) for an address
func (j *JSONRPCServiceImpl) GetTransactionCount(ctx context.Context, address common.Address, blockNumber *big.Int) (uint64, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_getTransactionCount"]++

	// TODO: Implement nonce lookup from state
	return 0, fmt.Errorf("not implemented")
}

// Account related methods

// GetBalance returns the balance of an address
func (j *JSONRPCServiceImpl) GetBalance(ctx context.Context, address common.Address, blockNumber *big.Int) (*big.Int, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_getBalance"]++

	// TODO: Implement balance lookup from state
	return big.NewInt(0), nil
}

// GetCode returns the contract code at an address
func (j *JSONRPCServiceImpl) GetCode(ctx context.Context, address common.Address, blockNumber *big.Int) (hexutil.Bytes, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_getCode"]++

	// TODO: Implement code lookup from state
	return hexutil.Bytes{}, nil
}

// GetStorageAt returns the storage value at a specific slot
func (j *JSONRPCServiceImpl) GetStorageAt(ctx context.Context, address common.Address, key common.Hash, blockNumber *big.Int) (common.Hash, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_getStorageAt"]++

	// TODO: Implement storage lookup from state
	return common.Hash{}, nil
}

// Block related methods

// GetBlockByNumber returns a block by its number
func (j *JSONRPCServiceImpl) GetBlockByNumber(ctx context.Context, number *big.Int, fullTx bool) (*types.Block, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_getBlockByNumber"]++

	// TODO: Implement block lookup from state
	return nil, fmt.Errorf("not implemented")
}

// GetBlockByHash returns a block by its hash
func (j *JSONRPCServiceImpl) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (*types.Block, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_getBlockByHash"]++

	// TODO: Implement block lookup from state
	return nil, fmt.Errorf("not implemented")
}

// BlockNumber returns the current block number
func (j *JSONRPCServiceImpl) BlockNumber(ctx context.Context) (uint64, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_blockNumber"]++

	// TODO: Implement current block number lookup
	return 0, nil
}

// Call methods

// Call executes a contract call without creating a transaction
func (j *JSONRPCServiceImpl) Call(ctx context.Context, args interfaces.CallArgs, blockNumber *big.Int) (hexutil.Bytes, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_call"]++

	// TODO: Implement call execution using EVM
	return hexutil.Bytes{}, nil
}

// EstimateGas estimates the gas required for a transaction
func (j *JSONRPCServiceImpl) EstimateGas(ctx context.Context, args interfaces.CallArgs, blockNumber *big.Int) (uint64, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_estimateGas"]++

	// TODO: Implement gas estimation
	return 21000, nil // Return minimum gas for now
}

// Network methods

// ChainId returns the chain ID
func (j *JSONRPCServiceImpl) ChainId(ctx context.Context) (*big.Int, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_chainId"]++

	// TODO: Return configured chain ID
	return big.NewInt(1337), nil
}

// GasPrice returns the current gas price
func (j *JSONRPCServiceImpl) GasPrice(ctx context.Context) (*big.Int, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_gasPrice"]++

	// TODO: Implement dynamic gas pricing
	return big.NewInt(20000000000), nil // 20 gwei
}

// Log methods

// GetLogs returns logs matching the filter criteria
func (j *JSONRPCServiceImpl) GetLogs(ctx context.Context, crit interfaces.FilterQuery) ([]*types.Log, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_getLogs"]++

	// TODO: Implement log filtering and retrieval
	return []*types.Log{}, nil
}

// Utility methods

// GetProof returns account and storage proofs
func (j *JSONRPCServiceImpl) GetProof(ctx context.Context, address common.Address, keys []common.Hash, blockNumber *big.Int) (*interfaces.AccountResult, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.metrics.TotalRequests++
	j.metrics.RequestsByMethod["eth_getProof"]++

	// TODO: Implement proof generation
	return nil, fmt.Errorf("not implemented")
}

// SetTxPoolService sets the transaction pool service
func (j *JSONRPCServiceImpl) SetTxPoolService(service interfaces.TxPoolService) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.txPoolService = service
}

// SetStateReader sets the state reader
func (j *JSONRPCServiceImpl) SetStateReader(reader sharedinterfaces.StateReader) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.stateReader = reader
}

// SetMessageProducer sets the message producer
func (j *JSONRPCServiceImpl) SetMessageProducer(producer sharedinterfaces.MessageProducer) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.messageProducer = producer
}

// SetRequestValidator sets the request validator
func (j *JSONRPCServiceImpl) SetRequestValidator(validator interfaces.RequestValidator) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.requestValidator = validator
}

// GetMetrics returns JSON-RPC service metrics
func (j *JSONRPCServiceImpl) GetMetrics() *interfaces.JSONRPCMetrics {
	j.mu.RLock()
	defer j.mu.RUnlock()

	// Calculate uptime and request rate
	uptime := time.Since(j.startTime)
	rps := float64(j.metrics.TotalRequests) / uptime.Seconds()

	metrics := *j.metrics
	metrics.RequestsPerSecond = rps
	metrics.Timestamp = time.Now()

	return &metrics
}

// HandleHTTPRequest handles HTTP JSON-RPC requests
func (j *JSONRPCServiceImpl) HandleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST and GET methods
	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var requestData json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		j.sendJSONRPCError(w, nil, -32700, "Parse error", nil)
		return
	}

	// Check if it's a batch request
	if len(requestData) > 0 && requestData[0] == '[' {
		j.handleBatchRequest(w, requestData)
	} else {
		j.handleSingleRequest(w, requestData)
	}
}

// handleSingleRequest handles a single JSON-RPC request
func (j *JSONRPCServiceImpl) handleSingleRequest(w http.ResponseWriter, requestData json.RawMessage) {
	var request interfaces.JSONRPCRequest
	if err := json.Unmarshal(requestData, &request); err != nil {
		j.sendJSONRPCError(w, nil, -32700, "Parse error", err.Error())
		return
	}

	// Validate JSON-RPC version
	if request.JSONRPC != "2.0" {
		j.sendJSONRPCError(w, request.ID, -32600, "Invalid Request", "jsonrpc must be 2.0")
		return
	}

	// Route and execute the method
	result, err := j.executeMethod(context.Background(), &request)
	if err != nil {
		j.sendJSONRPCError(w, request.ID, -32603, "Internal error", err.Error())
		return
	}

	// Send successful response
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      request.ID,
		"result":  result,
	}

	json.NewEncoder(w).Encode(response)
}

// handleBatchRequest handles a batch of JSON-RPC requests
func (j *JSONRPCServiceImpl) handleBatchRequest(w http.ResponseWriter, requestData json.RawMessage) {
	var requests []interfaces.JSONRPCRequest
	if err := json.Unmarshal(requestData, &requests); err != nil {
		j.sendJSONRPCError(w, nil, -32700, "Parse error", err.Error())
		return
	}

	// Check batch size limit
	if len(requests) > j.config.BatchLimit {
		j.sendJSONRPCError(w, nil, -32600, "Invalid Request", "Batch size exceeds limit")
		return
	}

	responses := make([]map[string]interface{}, 0, len(requests))

	for _, request := range requests {
		if request.JSONRPC != "2.0" {
			responses = append(responses, map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      request.ID,
				"error": map[string]interface{}{
					"code":    -32600,
					"message": "Invalid Request",
					"data":    "jsonrpc must be 2.0",
				},
			})
			continue
		}

		result, err := j.executeMethod(context.Background(), &request)
		if err != nil {
			responses = append(responses, map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      request.ID,
				"error": map[string]interface{}{
					"code":    -32603,
					"message": "Internal error",
					"data":    err.Error(),
				},
			})
		} else {
			responses = append(responses, map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      request.ID,
				"result":  result,
			})
		}
	}

	json.NewEncoder(w).Encode(responses)
}

// executeMethod executes a JSON-RPC method
func (j *JSONRPCServiceImpl) executeMethod(ctx context.Context, request *interfaces.JSONRPCRequest) (interface{}, error) {
	switch request.Method {
	// Transaction methods
	case "eth_sendRawTransaction":
		return j.executeSendRawTransaction(ctx, request.Params)
	case "eth_getTransactionByHash":
		return j.executeGetTransactionByHash(ctx, request.Params)
	case "eth_getTransactionReceipt":
		return j.executeGetTransactionReceipt(ctx, request.Params)
	case "eth_getTransactionCount":
		return j.executeGetTransactionCount(ctx, request.Params)

	// Account methods
	case "eth_getBalance":
		return j.executeGetBalance(ctx, request.Params)
	case "eth_getCode":
		return j.executeGetCode(ctx, request.Params)
	case "eth_getStorageAt":
		return j.executeGetStorageAt(ctx, request.Params)

	// Block methods
	case "eth_getBlockByNumber":
		return j.executeGetBlockByNumber(ctx, request.Params)
	case "eth_getBlockByHash":
		return j.executeGetBlockByHash(ctx, request.Params)
	case "eth_blockNumber":
		return j.executeBlockNumber(ctx, request.Params)

	// Call methods
	case "eth_call":
		return j.executeCall(ctx, request.Params)
	case "eth_estimateGas":
		return j.executeEstimateGas(ctx, request.Params)

	// Network methods
	case "eth_chainId":
		return j.executeChainId(ctx, request.Params)
	case "eth_gasPrice":
		return j.executeGasPrice(ctx, request.Params)
	case "net_version":
		return j.executeNetVersion(ctx, request.Params)

	// Log methods
	case "eth_getLogs":
		return j.executeGetLogs(ctx, request.Params)

	// Utility methods
	case "eth_getProof":
		return j.executeGetProof(ctx, request.Params)

	// Node information
	case "web3_clientVersion":
		return "YaoPortal/v1.0.0", nil
	case "net_listening":
		return true, nil
	case "net_peerCount":
		return "0x0", nil

	default:
		return nil, fmt.Errorf("method not found: %s", request.Method)
	}
}

// Execute method implementations
func (j *JSONRPCServiceImpl) executeSendRawTransaction(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) != 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	dataHex, ok := paramsList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid data parameter")
	}

	data, err := hexutil.Decode(dataHex)
	if err != nil {
		return nil, fmt.Errorf("invalid hex data: %w", err)
	}

	hash, err := j.SendRawTransaction(ctx, data)
	return hash.Hex(), err
}

func (j *JSONRPCServiceImpl) executeGetTransactionByHash(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) != 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	hashHex, ok := paramsList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid hash parameter")
	}

	hash := common.HexToHash(hashHex)
	return j.GetTransactionByHash(ctx, hash)
}

func (j *JSONRPCServiceImpl) executeGetTransactionReceipt(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) != 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	hashHex, ok := paramsList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid hash parameter")
	}

	hash := common.HexToHash(hashHex)
	return j.GetTransactionReceipt(ctx, hash)
}

func (j *JSONRPCServiceImpl) executeGetTransactionCount(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	addressHex, ok := paramsList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid address parameter")
	}

	address := common.HexToAddress(addressHex)
	var blockNumber *big.Int
	if len(paramsList) > 1 {
		blockNumStr, ok := paramsList[1].(string)
		if ok && blockNumStr != "latest" {
			blockNumber, _ = hexutil.DecodeBig(blockNumStr)
		}
	}

	count, err := j.GetTransactionCount(ctx, address, blockNumber)
	return fmt.Sprintf("0x%x", count), err
}

func (j *JSONRPCServiceImpl) executeGetBalance(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	addressHex, ok := paramsList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid address parameter")
	}

	address := common.HexToAddress(addressHex)
	var blockNumber *big.Int
	if len(paramsList) > 1 {
		blockNumStr, ok := paramsList[1].(string)
		if ok && blockNumStr != "latest" {
			blockNumber, _ = hexutil.DecodeBig(blockNumStr)
		}
	}

	balance, err := j.GetBalance(ctx, address, blockNumber)
	return (*hexutil.Big)(balance), err
}

func (j *JSONRPCServiceImpl) executeGetCode(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	addressHex, ok := paramsList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid address parameter")
	}

	address := common.HexToAddress(addressHex)
	var blockNumber *big.Int
	if len(paramsList) > 1 {
		blockNumStr, ok := paramsList[1].(string)
		if ok && blockNumStr != "latest" {
			blockNumber, _ = hexutil.DecodeBig(blockNumStr)
		}
	}

	return j.GetCode(ctx, address, blockNumber)
}

func (j *JSONRPCServiceImpl) executeGetStorageAt(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) < 2 {
		return nil, fmt.Errorf("invalid parameters")
	}

	addressHex, ok := paramsList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid address parameter")
	}

	keyHex, ok := paramsList[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid key parameter")
	}

	address := common.HexToAddress(addressHex)
	key := common.HexToHash(keyHex)
	var blockNumber *big.Int
	if len(paramsList) > 2 {
		blockNumStr, ok := paramsList[2].(string)
		if ok && blockNumStr != "latest" {
			blockNumber, _ = hexutil.DecodeBig(blockNumStr)
		}
	}

	result, err := j.GetStorageAt(ctx, address, key, blockNumber)
	return result.Hex(), err
}

func (j *JSONRPCServiceImpl) executeGetBlockByNumber(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) < 2 {
		return nil, fmt.Errorf("invalid parameters")
	}

	blockNumParam, ok := paramsList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid block number parameter")
	}

	fullTx, ok := paramsList[1].(bool)
	if !ok {
		return nil, fmt.Errorf("invalid fullTx parameter")
	}

	var blockNumber *big.Int
	if blockNumParam != "latest" && blockNumParam != "earliest" && blockNumParam != "pending" {
		blockNumber, _ = hexutil.DecodeBig(blockNumParam)
	}

	return j.GetBlockByNumber(ctx, blockNumber, fullTx)
}

func (j *JSONRPCServiceImpl) executeGetBlockByHash(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) < 2 {
		return nil, fmt.Errorf("invalid parameters")
	}

	hashHex, ok := paramsList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid hash parameter")
	}

	fullTx, ok := paramsList[1].(bool)
	if !ok {
		return nil, fmt.Errorf("invalid fullTx parameter")
	}

	hash := common.HexToHash(hashHex)
	return j.GetBlockByHash(ctx, hash, fullTx)
}

func (j *JSONRPCServiceImpl) executeBlockNumber(ctx context.Context, params interface{}) (interface{}, error) {
	blockNum, err := j.BlockNumber(ctx)
	return fmt.Sprintf("0x%x", blockNum), err
}

func (j *JSONRPCServiceImpl) executeCall(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	callArgsMap, ok := paramsList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid call arguments")
	}

	args := interfaces.CallArgs{}
	if to, ok := callArgsMap["to"].(string); ok {
		addr := common.HexToAddress(to)
		args.To = &addr
	}
	if from, ok := callArgsMap["from"].(string); ok {
		addr := common.HexToAddress(from)
		args.From = &addr
	}
	if data, ok := callArgsMap["data"].(string); ok {
		args.Data, _ = hexutil.Decode(data)
	}

	var blockNumber *big.Int
	if len(paramsList) > 1 {
		blockNumStr, ok := paramsList[1].(string)
		if ok && blockNumStr != "latest" {
			blockNumber, _ = hexutil.DecodeBig(blockNumStr)
		}
	}

	return j.Call(ctx, args, blockNumber)
}

func (j *JSONRPCServiceImpl) executeEstimateGas(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	callArgsMap, ok := paramsList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid call arguments")
	}

	args := interfaces.CallArgs{}
	if to, ok := callArgsMap["to"].(string); ok {
		addr := common.HexToAddress(to)
		args.To = &addr
	}
	if from, ok := callArgsMap["from"].(string); ok {
		addr := common.HexToAddress(from)
		args.From = &addr
	}
	if data, ok := callArgsMap["data"].(string); ok {
		args.Data, _ = hexutil.Decode(data)
	}

	gas, err := j.EstimateGas(ctx, args, nil)
	return fmt.Sprintf("0x%x", gas), err
}

func (j *JSONRPCServiceImpl) executeChainId(ctx context.Context, params interface{}) (interface{}, error) {
	chainId, err := j.ChainId(ctx)
	return fmt.Sprintf("0x%x", chainId), err
}

func (j *JSONRPCServiceImpl) executeGasPrice(ctx context.Context, params interface{}) (interface{}, error) {
	gasPrice, err := j.GasPrice(ctx)
	return (*hexutil.Big)(gasPrice), err
}

func (j *JSONRPCServiceImpl) executeNetVersion(ctx context.Context, params interface{}) (interface{}, error) {
	chainId, err := j.ChainId(ctx)
	return chainId.String(), err
}

func (j *JSONRPCServiceImpl) executeGetLogs(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	_, ok = paramsList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid filter criteria")
	}

	filter := interfaces.FilterQuery{}
	// Parse filter parameters
	// TODO: Implement full filter parsing

	return j.GetLogs(ctx, filter)
}

func (j *JSONRPCServiceImpl) executeGetProof(ctx context.Context, params interface{}) (interface{}, error) {
	paramsList, ok := params.([]interface{})
	if !ok || len(paramsList) < 2 {
		return nil, fmt.Errorf("invalid parameters")
	}

	addressHex, ok := paramsList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid address parameter")
	}

	keysParam, ok := paramsList[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid keys parameter")
	}

	address := common.HexToAddress(addressHex)
	keys := make([]common.Hash, len(keysParam))
	for i, keyParam := range keysParam {
		keyHex, ok := keyParam.(string)
		if !ok {
			return nil, fmt.Errorf("invalid key parameter")
		}
		keys[i] = common.HexToHash(keyHex)
	}

	var blockNumber *big.Int
	if len(paramsList) > 2 {
		blockNumStr, ok := paramsList[2].(string)
		if ok && blockNumStr != "latest" {
			blockNumber, _ = hexutil.DecodeBig(blockNumStr)
		}
	}

	return j.GetProof(ctx, address, keys, blockNumber)
}

// sendJSONRPCError sends a JSON-RPC error response
func (j *JSONRPCServiceImpl) sendJSONRPCError(w http.ResponseWriter, id interface{}, code int, message string, data interface{}) {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"data":    data,
		},
	}

	w.WriteHeader(http.StatusOK) // JSON-RPC errors still return 200
	json.NewEncoder(w).Encode(response)
}
