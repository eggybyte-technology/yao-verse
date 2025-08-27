// Package types contains blockchain-related types for YaoVerse
package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	Hash     common.Hash     `json:"hash"`
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to,omitempty"`
	Value    *big.Int        `json:"value"`
	Gas      uint64          `json:"gas"`
	GasPrice *big.Int        `json:"gasPrice"`
	Data     []byte          `json:"data"`
	Nonce    uint64          `json:"nonce"`
	// EIP-2930
	AccessList []AccessTuple `json:"accessList,omitempty"`
	// EIP-1559
	MaxFeePerGas         *big.Int `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas,omitempty"`
	// Signature fields
	V *big.Int `json:"v"`
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

// AccessTuple represents an access tuple for EIP-2930
type AccessTuple struct {
	Address     common.Address `json:"address"`
	StorageKeys []common.Hash  `json:"storageKeys"`
}

// Block represents a simplified block structure in YaoVerse
type Block struct {
	Number       uint64         `json:"number"`
	Hash         common.Hash    `json:"hash"`
	ParentHash   common.Hash    `json:"parentHash"`
	Timestamp    uint64         `json:"timestamp"`
	GasLimit     uint64         `json:"gasLimit"`
	GasUsed      uint64         `json:"gasUsed"`
	Transactions []*Transaction `json:"transactions"`
	ReceiptHash  common.Hash    `json:"receiptHash"`
	StateRoot    common.Hash    `json:"stateRoot"`
}

// Receipt represents a transaction receipt
type Receipt struct {
	TransactionHash   common.Hash     `json:"transactionHash"`
	TransactionIndex  uint64          `json:"transactionIndex"`
	BlockHash         common.Hash     `json:"blockHash"`
	BlockNumber       uint64          `json:"blockNumber"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to,omitempty"`
	GasUsed           uint64          `json:"gasUsed"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	ContractAddress   *common.Address `json:"contractAddress,omitempty"`
	Logs              []*Log          `json:"logs"`
	Status            uint64          `json:"status"`
	EffectiveGasPrice *big.Int        `json:"effectiveGasPrice"`
}

// Log represents an event log
type Log struct {
	Address          common.Address `json:"address"`
	Topics           []common.Hash  `json:"topics"`
	Data             []byte         `json:"data"`
	BlockNumber      uint64         `json:"blockNumber"`
	TransactionHash  common.Hash    `json:"transactionHash"`
	TransactionIndex uint64         `json:"transactionIndex"`
	BlockHash        common.Hash    `json:"blockHash"`
	LogIndex         uint64         `json:"logIndex"`
	Removed          bool           `json:"removed"`
}

// Account represents an account state
type Account struct {
	Address common.Address `json:"address"`
	Nonce   uint64         `json:"nonce"`
	Balance *big.Int       `json:"balance"`
	// Code hash for contract accounts
	CodeHash common.Hash `json:"codeHash"`
	// Storage root for contract accounts
	StorageRoot common.Hash `json:"storageRoot"`
}

// StorageEntry represents a single storage slot
type StorageEntry struct {
	Address common.Address `json:"address"`
	Key     common.Hash    `json:"key"`
	Value   common.Hash    `json:"value"`
}

// ExecutionResult represents the result of block execution
type ExecutionResult struct {
	BlockNumber  uint64      `json:"blockNumber"`
	BlockHash    common.Hash `json:"blockHash"`
	StateRoot    common.Hash `json:"stateRoot"`
	Receipts     []*Receipt  `json:"receipts"`
	StateDiff    *StateDiff  `json:"stateDiff"`
	GasUsed      uint64      `json:"gasUsed"`
	ExecutedAt   time.Time   `json:"executedAt"`
	ExecutorNode string      `json:"executorNode"`
}

// StateDiff represents the state changes after block execution
type StateDiff struct {
	ModifiedAccounts map[common.Address]*Account                    `json:"modifiedAccounts"`
	ModifiedStorage  map[common.Address]map[common.Hash]common.Hash `json:"modifiedStorage"`
	CreatedContracts map[common.Address][]byte                      `json:"createdContracts"`
	DeletedAccounts  []common.Address                               `json:"deletedAccounts"`
}

// CacheKey represents a cache key for state data
type CacheKey struct {
	Type    CacheKeyType `json:"type"`
	Address string       `json:"address"`
	Key     string       `json:"key,omitempty"`
}

// CacheKeyType defines different types of cache keys
type CacheKeyType string

const (
	CacheKeyTypeAccount CacheKeyType = "account"
	CacheKeyTypeStorage CacheKeyType = "storage"
	CacheKeyTypeCode    CacheKeyType = "code"
)

// CacheInvalidationMessage represents the message sent to invalidate caches
type CacheInvalidationMessage struct {
	BlockNumber uint64      `json:"blockNumber"`
	Keys        []*CacheKey `json:"keys"`
	Timestamp   time.Time   `json:"timestamp"`
}

// ErrorCode defines error codes for YaoVerse
type ErrorCode int

const (
	ErrCodeSuccess ErrorCode = iota
	ErrCodeInvalidTransaction
	ErrCodeInsufficientBalance
	ErrCodeContractExecutionFailed
	ErrCodeStateNotFound
	ErrCodeInternalError
	ErrCodeRPCError
	ErrCodeConsensusError
	ErrCodeStorageError
)

// YaoError represents a YaoVerse specific error
type YaoError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Data    any       `json:"data,omitempty"`
}

func (e *YaoError) Error() string {
	return e.Message
}

// NewYaoError creates a new YaoError
func NewYaoError(code ErrorCode, message string) *YaoError {
	return &YaoError{
		Code:    code,
		Message: message,
	}
}
