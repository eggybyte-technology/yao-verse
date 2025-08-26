package storage

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/eggybyte-technology/yao-verse-shared/interfaces"
	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/go-sql-driver/mysql"
)

// TiDBStorage implements L3Storage interface using TiDB
type TiDBStorage struct {
	db     *sql.DB
	config *TiDBConfig
}

// TiDBConfig defines configuration for TiDB storage
type TiDBConfig struct {
	DSN             string `json:"dsn"`
	MaxOpenConns    int    `json:"maxOpenConns"`
	MaxIdleConns    int    `json:"maxIdleConns"`
	ConnMaxLifetime int    `json:"connMaxLifetime"` // in seconds
}

// NewTiDBStorage creates a new TiDB storage instance
func NewTiDBStorage(config *TiDBConfig) (*TiDBStorage, error) {
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open TiDB connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping TiDB: %w", err)
	}

	storage := &TiDBStorage{
		db:     db,
		config: config,
	}

	// Initialize database schema
	if err := storage.initializeSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return storage, nil
}

// GetAccount retrieves an account by address
func (s *TiDBStorage) GetAccount(ctx context.Context, address common.Address) (*types.Account, error) {
	query := `SELECT address, nonce, balance, code_hash, storage_root FROM accounts WHERE address = ?`

	row := s.db.QueryRowContext(ctx, query, address.Hex())

	var account types.Account
	var balanceStr, codeHashStr, storageRootStr string

	err := row.Scan(&account.Address, &account.Nonce, &balanceStr, &codeHashStr, &storageRootStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found: %s", address.Hex())
		}
		return nil, fmt.Errorf("failed to scan account: %w", err)
	}

	// Parse balance
	balance, ok := new(big.Int).SetString(balanceStr, 10)
	if !ok {
		return nil, fmt.Errorf("invalid balance format: %s", balanceStr)
	}
	account.Balance = balance
	account.CodeHash = common.HexToHash(codeHashStr)
	account.StorageRoot = common.HexToHash(storageRootStr)

	return &account, nil
}

// SetAccount sets account data
func (s *TiDBStorage) SetAccount(ctx context.Context, address common.Address, account *types.Account) error {
	query := `INSERT INTO accounts (address, nonce, balance, code_hash, storage_root, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, NOW(), NOW())
			  ON DUPLICATE KEY UPDATE 
			  nonce = VALUES(nonce), 
			  balance = VALUES(balance), 
			  code_hash = VALUES(code_hash), 
			  storage_root = VALUES(storage_root), 
			  updated_at = NOW()`

	_, err := s.db.ExecContext(ctx, query,
		address.Hex(),
		account.Nonce,
		account.Balance.String(),
		account.CodeHash.Hex(),
		account.StorageRoot.Hex())

	if err != nil {
		return fmt.Errorf("failed to set account: %w", err)
	}

	return nil
}

// DeleteAccount deletes an account
func (s *TiDBStorage) DeleteAccount(ctx context.Context, address common.Address) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete contract storage if it exists
	_, err = tx.ExecContext(ctx, "DELETE FROM contract_storage WHERE address = ?", address.Hex())
	if err != nil {
		return fmt.Errorf("failed to delete contract storage: %w", err)
	}

	// Delete contract code if it exists
	_, err = tx.ExecContext(ctx, "DELETE FROM contract_code WHERE address = ?", address.Hex())
	if err != nil {
		return fmt.Errorf("failed to delete contract code: %w", err)
	}

	// Delete account
	_, err = tx.ExecContext(ctx, "DELETE FROM accounts WHERE address = ?", address.Hex())
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return tx.Commit()
}

// HasAccount checks if an account exists
func (s *TiDBStorage) HasAccount(ctx context.Context, address common.Address) (bool, error) {
	query := `SELECT 1 FROM accounts WHERE address = ? LIMIT 1`

	var exists int
	err := s.db.QueryRowContext(ctx, query, address.Hex()).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check account existence: %w", err)
	}

	return true, nil
}

// GetStorageAt retrieves a storage value at a specific key
func (s *TiDBStorage) GetStorageAt(ctx context.Context, address common.Address, key common.Hash) (common.Hash, error) {
	query := `SELECT value FROM contract_storage WHERE address = ? AND key = ?`

	var valueStr string
	err := s.db.QueryRowContext(ctx, query, address.Hex(), key.Hex()).Scan(&valueStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return common.Hash{}, nil // Return empty hash for non-existent keys
		}
		return common.Hash{}, fmt.Errorf("failed to get storage: %w", err)
	}

	return common.HexToHash(valueStr), nil
}

// SetStorageAt sets a storage value at a specific key
func (s *TiDBStorage) SetStorageAt(ctx context.Context, address common.Address, key, value common.Hash) error {
	query := `INSERT INTO contract_storage (address, key, value, created_at, updated_at) 
			  VALUES (?, ?, ?, NOW(), NOW())
			  ON DUPLICATE KEY UPDATE 
			  value = VALUES(value), 
			  updated_at = NOW()`

	_, err := s.db.ExecContext(ctx, query, address.Hex(), key.Hex(), value.Hex())
	if err != nil {
		return fmt.Errorf("failed to set storage: %w", err)
	}

	return nil
}

// GetCode retrieves the code of a contract
func (s *TiDBStorage) GetCode(ctx context.Context, address common.Address) ([]byte, error) {
	query := `SELECT code FROM contract_code WHERE address = ?`

	var codeHex string
	err := s.db.QueryRowContext(ctx, query, address.Hex()).Scan(&codeHex)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil for accounts without code
		}
		return nil, fmt.Errorf("failed to get code: %w", err)
	}

	return common.FromHex(codeHex), nil
}

// SetCode sets the code of a contract
func (s *TiDBStorage) SetCode(ctx context.Context, address common.Address, code []byte) error {
	query := `INSERT INTO contract_code (address, code, code_hash, created_at, updated_at) 
			  VALUES (?, ?, ?, NOW(), NOW())
			  ON DUPLICATE KEY UPDATE 
			  code = VALUES(code), 
			  code_hash = VALUES(code_hash), 
			  updated_at = NOW()`

	codeHash := crypto.Keccak256Hash(code)
	_, err := s.db.ExecContext(ctx, query, address.Hex(), common.Bytes2Hex(code), codeHash.Hex())
	if err != nil {
		return fmt.Errorf("failed to set code: %w", err)
	}

	return nil
}

// GetCodeHash retrieves the code hash of a contract
func (s *TiDBStorage) GetCodeHash(ctx context.Context, address common.Address) (common.Hash, error) {
	query := `SELECT code_hash FROM contract_code WHERE address = ?`

	var codeHashStr string
	err := s.db.QueryRowContext(ctx, query, address.Hex()).Scan(&codeHashStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return common.Hash{}, nil
		}
		return common.Hash{}, fmt.Errorf("failed to get code hash: %w", err)
	}

	return common.HexToHash(codeHashStr), nil
}

// GetBlock retrieves a block by number
func (s *TiDBStorage) GetBlock(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	query := `SELECT block_number, block_hash, parent_hash, timestamp, gas_limit, gas_used, receipt_hash, state_root 
			  FROM blocks WHERE block_number = ?`

	row := s.db.QueryRowContext(ctx, query, blockNumber)

	var block types.Block
	var blockHashStr, parentHashStr, receiptHashStr, stateRootStr string

	err := row.Scan(&block.Number, &blockHashStr, &parentHashStr, &block.Timestamp,
		&block.GasLimit, &block.GasUsed, &receiptHashStr, &stateRootStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("block not found: %d", blockNumber)
		}
		return nil, fmt.Errorf("failed to scan block: %w", err)
	}

	block.Hash = common.HexToHash(blockHashStr)
	block.ParentHash = common.HexToHash(parentHashStr)
	block.ReceiptHash = common.HexToHash(receiptHashStr)
	block.StateRoot = common.HexToHash(stateRootStr)

	// Load transactions for this block
	transactions, err := s.getBlockTransactions(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to load block transactions: %w", err)
	}
	block.Transactions = transactions

	return &block, nil
}

// GetBlockByHash retrieves a block by hash
func (s *TiDBStorage) GetBlockByHash(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	query := `SELECT block_number FROM blocks WHERE block_hash = ?`

	var blockNumber uint64
	err := s.db.QueryRowContext(ctx, query, blockHash.Hex()).Scan(&blockNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("block not found: %s", blockHash.Hex())
		}
		return nil, fmt.Errorf("failed to get block number: %w", err)
	}

	return s.GetBlock(ctx, blockNumber)
}

// SetBlock sets block data
func (s *TiDBStorage) SetBlock(ctx context.Context, block *types.Block) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert block
	blockQuery := `INSERT INTO blocks (block_number, block_hash, parent_hash, timestamp, gas_limit, gas_used, 
						receipt_hash, state_root, created_at) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())`

	_, err = tx.ExecContext(ctx, blockQuery, block.Number, block.Hash.Hex(), block.ParentHash.Hex(),
		block.Timestamp, block.GasLimit, block.GasUsed, block.ReceiptHash.Hex(),
		block.StateRoot.Hex())
	if err != nil {
		return fmt.Errorf("failed to insert block: %w", err)
	}

	// Insert transactions
	for i, transaction := range block.Transactions {
		err = s.insertTransaction(ctx, transaction, block.Number, uint64(i), tx)
		if err != nil {
			return fmt.Errorf("failed to insert transaction %s: %w", transaction.Hash.Hex(), err)
		}
	}

	return tx.Commit()
}

// BeginTransaction begins a database transaction
func (s *TiDBStorage) BeginTransaction(ctx context.Context) (interfaces.Transaction, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &tidbTransaction{tx: tx, storage: s}, nil
}

// GetStateRoot returns the state root for a given block
func (s *TiDBStorage) GetStateRoot(ctx context.Context, blockNumber uint64) (common.Hash, error) {
	query := `SELECT state_root FROM blocks WHERE block_number = ?`

	var stateRootStr string
	err := s.db.QueryRowContext(ctx, query, blockNumber).Scan(&stateRootStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return common.Hash{}, fmt.Errorf("block not found: %d", blockNumber)
		}
		return common.Hash{}, fmt.Errorf("failed to get state root: %w", err)
	}

	return common.HexToHash(stateRootStr), nil
}

// SetStateRoot sets the state root for a given block
func (s *TiDBStorage) SetStateRoot(ctx context.Context, blockNumber uint64, stateRoot common.Hash) error {
	query := `UPDATE blocks SET state_root = ?, updated_at = NOW() WHERE block_number = ?`

	result, err := s.db.ExecContext(ctx, query, stateRoot.Hex(), blockNumber)
	if err != nil {
		return fmt.Errorf("failed to set state root: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("block not found: %d", blockNumber)
	}

	return nil
}

// GetTransaction retrieves a transaction by hash
func (s *TiDBStorage) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, error) {
	query := `SELECT tx_hash, from_address, to_address, value, gas, gas_price, data, nonce, 
			  max_fee_per_gas, max_priority_fee_per_gas, v, r, s 
			  FROM transactions WHERE tx_hash = ?`

	row := s.db.QueryRowContext(ctx, query, txHash.Hex())

	var tx types.Transaction
	var fromStr, valueStr, gasPriceStr, dataHex, vStr, rStr, sStr string
	var toAddress sql.NullString
	var maxFee, maxPriority sql.NullString

	err := row.Scan(&tx.Hash, &fromStr, &toAddress, &valueStr, &tx.Gas, &gasPriceStr, &dataHex,
		&tx.Nonce, &maxFee, &maxPriority, &vStr, &rStr, &sStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found: %s", txHash.Hex())
		}
		return nil, fmt.Errorf("failed to scan transaction: %w", err)
	}

	// Parse fields
	tx.From = common.HexToAddress(fromStr)
	if toAddress.Valid {
		addr := common.HexToAddress(toAddress.String)
		tx.To = &addr
	}

	value, ok := new(big.Int).SetString(valueStr, 10)
	if !ok {
		return nil, fmt.Errorf("invalid value format: %s", valueStr)
	}
	tx.Value = value

	gasPrice, ok := new(big.Int).SetString(gasPriceStr, 10)
	if !ok {
		return nil, fmt.Errorf("invalid gas price format: %s", gasPriceStr)
	}
	tx.GasPrice = gasPrice

	tx.Data = common.FromHex(dataHex)

	// Parse EIP-1559 fields
	if maxFee.Valid {
		maxFeePerGas, ok := new(big.Int).SetString(maxFee.String, 10)
		if ok {
			tx.MaxFeePerGas = maxFeePerGas
		}
	}

	if maxPriority.Valid {
		maxPriorityPerGas, ok := new(big.Int).SetString(maxPriority.String, 10)
		if ok {
			tx.MaxPriorityFeePerGas = maxPriorityPerGas
		}
	}

	// Parse signature
	tx.V, _ = new(big.Int).SetString(vStr, 10)
	tx.R, _ = new(big.Int).SetString(rStr, 10)
	tx.S, _ = new(big.Int).SetString(sStr, 10)

	return &tx, nil
}

// GetTransactionReceipt retrieves a transaction receipt by hash
func (s *TiDBStorage) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	query := `SELECT tx_hash, tx_index, block_hash, block_number, from_address, to_address, 
			  gas_used, cumulative_gas_used, contract_address, status, effective_gas_price
			  FROM transaction_receipts WHERE tx_hash = ?`

	row := s.db.QueryRowContext(ctx, query, txHash.Hex())

	var receipt types.Receipt
	var blockHashStr, fromStr, effectiveGasPriceStr string
	var toAddress, contractAddress sql.NullString

	err := row.Scan(&receipt.TransactionHash, &receipt.TransactionIndex, &blockHashStr,
		&receipt.BlockNumber, &fromStr, &toAddress, &receipt.GasUsed,
		&receipt.CumulativeGasUsed, &contractAddress, &receipt.Status, &effectiveGasPriceStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("receipt not found: %s", txHash.Hex())
		}
		return nil, fmt.Errorf("failed to scan receipt: %w", err)
	}

	receipt.BlockHash = common.HexToHash(blockHashStr)
	receipt.From = common.HexToAddress(fromStr)

	if toAddress.Valid {
		addr := common.HexToAddress(toAddress.String)
		receipt.To = &addr
	}

	if contractAddress.Valid {
		addr := common.HexToAddress(contractAddress.String)
		receipt.ContractAddress = &addr
	}

	effectiveGasPrice, ok := new(big.Int).SetString(effectiveGasPriceStr, 10)
	if ok {
		receipt.EffectiveGasPrice = effectiveGasPrice
	}

	// Load logs for this receipt
	logs, err := s.getReceiptLogs(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to load receipt logs: %w", err)
	}
	receipt.Logs = logs

	return &receipt, nil
}

// Close closes the database connection
func (s *TiDBStorage) Close() error {
	return s.db.Close()
}

// Helper methods

// initializeSchema initializes the database schema
func (s *TiDBStorage) initializeSchema() error {
	schemas := []string{
		// Accounts table
		`CREATE TABLE IF NOT EXISTS accounts (
			address VARCHAR(42) PRIMARY KEY,
			nonce BIGINT UNSIGNED NOT NULL DEFAULT 0,
			balance TEXT NOT NULL DEFAULT '0',
			code_hash VARCHAR(66) NOT NULL DEFAULT '0x',
			storage_root VARCHAR(66) NOT NULL DEFAULT '0x',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_accounts_updated_at (updated_at)
		) ENGINE=InnoDB`,

		// Contract storage table
		`CREATE TABLE IF NOT EXISTS contract_storage (
			address VARCHAR(42) NOT NULL,
			key VARCHAR(66) NOT NULL,
			value VARCHAR(66) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (address, key),
			INDEX idx_storage_address (address),
			INDEX idx_storage_updated_at (updated_at)
		) ENGINE=InnoDB`,

		// Contract code table
		`CREATE TABLE IF NOT EXISTS contract_code (
			address VARCHAR(42) PRIMARY KEY,
			code LONGTEXT NOT NULL,
			code_hash VARCHAR(66) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_code_hash (code_hash),
			INDEX idx_code_updated_at (updated_at)
		) ENGINE=InnoDB`,

		// Blocks table
		`CREATE TABLE IF NOT EXISTS blocks (
			block_number BIGINT UNSIGNED PRIMARY KEY,
			block_hash VARCHAR(66) UNIQUE NOT NULL,
			parent_hash VARCHAR(66) NOT NULL,
			timestamp BIGINT UNSIGNED NOT NULL,
			gas_limit BIGINT UNSIGNED NOT NULL,
			gas_used BIGINT UNSIGNED NOT NULL,
			receipt_hash VARCHAR(66) NOT NULL,
			state_root VARCHAR(66) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_blocks_hash (block_hash),
			INDEX idx_blocks_timestamp (timestamp),
			INDEX idx_blocks_created_at (created_at)
		) ENGINE=InnoDB`,

		// Transactions table
		`CREATE TABLE IF NOT EXISTS transactions (
			tx_hash VARCHAR(66) PRIMARY KEY,
			block_number BIGINT UNSIGNED NOT NULL,
			tx_index INT UNSIGNED NOT NULL,
			from_address VARCHAR(42) NOT NULL,
			to_address VARCHAR(42) DEFAULT NULL,
			value TEXT NOT NULL DEFAULT '0',
			gas BIGINT UNSIGNED NOT NULL,
			gas_price TEXT NOT NULL DEFAULT '0',
			data LONGTEXT DEFAULT NULL,
			nonce BIGINT UNSIGNED NOT NULL,
			max_fee_per_gas TEXT DEFAULT NULL,
			max_priority_fee_per_gas TEXT DEFAULT NULL,
			v TEXT NOT NULL,
			r TEXT NOT NULL,
			s TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_tx_block (block_number),
			INDEX idx_tx_from (from_address),
			INDEX idx_tx_to (to_address),
			INDEX idx_tx_created_at (created_at),
			FOREIGN KEY (block_number) REFERENCES blocks(block_number)
		) ENGINE=InnoDB`,

		// Transaction receipts table
		`CREATE TABLE IF NOT EXISTS transaction_receipts (
			tx_hash VARCHAR(66) PRIMARY KEY,
			tx_index INT UNSIGNED NOT NULL,
			block_hash VARCHAR(66) NOT NULL,
			block_number BIGINT UNSIGNED NOT NULL,
			from_address VARCHAR(42) NOT NULL,
			to_address VARCHAR(42) DEFAULT NULL,
			gas_used BIGINT UNSIGNED NOT NULL,
			cumulative_gas_used BIGINT UNSIGNED NOT NULL,
			contract_address VARCHAR(42) DEFAULT NULL,
			status TINYINT UNSIGNED NOT NULL,
			effective_gas_price TEXT NOT NULL DEFAULT '0',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_receipt_block (block_number),
			INDEX idx_receipt_from (from_address),
			INDEX idx_receipt_to (to_address),
			INDEX idx_receipt_contract (contract_address),
			INDEX idx_receipt_created_at (created_at),
			FOREIGN KEY (block_number) REFERENCES blocks(block_number)
		) ENGINE=InnoDB`,

		// Event logs table
		`CREATE TABLE IF NOT EXISTS event_logs (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			address VARCHAR(42) NOT NULL,
			topic0 VARCHAR(66) DEFAULT NULL,
			topic1 VARCHAR(66) DEFAULT NULL,
			topic2 VARCHAR(66) DEFAULT NULL,
			topic3 VARCHAR(66) DEFAULT NULL,
			data LONGTEXT DEFAULT NULL,
			block_number BIGINT UNSIGNED NOT NULL,
			tx_hash VARCHAR(66) NOT NULL,
			tx_index INT UNSIGNED NOT NULL,
			block_hash VARCHAR(66) NOT NULL,
			log_index INT UNSIGNED NOT NULL,
			removed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_logs_address (address),
			INDEX idx_logs_topic0 (topic0),
			INDEX idx_logs_topic1 (topic1),
			INDEX idx_logs_topic2 (topic2),
			INDEX idx_logs_topic3 (topic3),
			INDEX idx_logs_block (block_number),
			INDEX idx_logs_tx (tx_hash),
			INDEX idx_logs_created_at (created_at),
			FOREIGN KEY (block_number) REFERENCES blocks(block_number),
			FOREIGN KEY (tx_hash) REFERENCES transactions(tx_hash)
		) ENGINE=InnoDB`,
	}

	for _, schema := range schemas {
		_, err := s.db.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// getBlockTransactions retrieves all transactions for a block
func (s *TiDBStorage) getBlockTransactions(ctx context.Context, blockNumber uint64) ([]*types.Transaction, error) {
	query := `SELECT tx_hash, from_address, to_address, value, gas, gas_price, data, nonce,
			  max_fee_per_gas, max_priority_fee_per_gas, v, r, s 
			  FROM transactions WHERE block_number = ? ORDER BY tx_index`

	rows, err := s.db.QueryContext(ctx, query, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*types.Transaction
	for rows.Next() {
		var tx types.Transaction
		var fromStr, valueStr, gasPriceStr, dataHex, vStr, rStr, sStr string
		var toAddress sql.NullString
		var maxFee, maxPriority sql.NullString

		err := rows.Scan(&tx.Hash, &fromStr, &toAddress, &valueStr, &tx.Gas, &gasPriceStr,
			&dataHex, &tx.Nonce, &maxFee, &maxPriority, &vStr, &rStr, &sStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Parse fields (similar to GetTransaction)
		tx.From = common.HexToAddress(fromStr)
		if toAddress.Valid {
			addr := common.HexToAddress(toAddress.String)
			tx.To = &addr
		}

		tx.Value, _ = new(big.Int).SetString(valueStr, 10)
		tx.GasPrice, _ = new(big.Int).SetString(gasPriceStr, 10)
		tx.Data = common.FromHex(dataHex)

		if maxFee.Valid {
			tx.MaxFeePerGas, _ = new(big.Int).SetString(maxFee.String, 10)
		}
		if maxPriority.Valid {
			tx.MaxPriorityFeePerGas, _ = new(big.Int).SetString(maxPriority.String, 10)
		}

		tx.V, _ = new(big.Int).SetString(vStr, 10)
		tx.R, _ = new(big.Int).SetString(rStr, 10)
		tx.S, _ = new(big.Int).SetString(sStr, 10)

		transactions = append(transactions, &tx)
	}

	return transactions, nil
}

// insertTransaction inserts a transaction
func (s *TiDBStorage) insertTransaction(ctx context.Context, tx *types.Transaction,
	blockNumber uint64, txIndex uint64, sqlTx *sql.Tx) error {
	query := `INSERT INTO transactions (tx_hash, block_number, tx_index, from_address, to_address, 
			  value, gas, gas_price, data, nonce, max_fee_per_gas, max_priority_fee_per_gas, 
			  v, r, s, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`

	var toAddr sql.NullString
	if tx.To != nil {
		toAddr.String = tx.To.Hex()
		toAddr.Valid = true
	}

	var maxFee, maxPriority sql.NullString
	if tx.MaxFeePerGas != nil {
		maxFee.String = tx.MaxFeePerGas.String()
		maxFee.Valid = true
	}
	if tx.MaxPriorityFeePerGas != nil {
		maxPriority.String = tx.MaxPriorityFeePerGas.String()
		maxPriority.Valid = true
	}

	_, err := sqlTx.ExecContext(ctx, query, tx.Hash.Hex(), blockNumber, txIndex,
		tx.From.Hex(), toAddr, tx.Value.String(), tx.Gas,
		tx.GasPrice.String(), common.Bytes2Hex(tx.Data), tx.Nonce,
		maxFee, maxPriority, tx.V.String(), tx.R.String(), tx.S.String())

	return err
}

// getReceiptLogs retrieves all logs for a receipt
func (s *TiDBStorage) getReceiptLogs(ctx context.Context, txHash common.Hash) ([]*types.Log, error) {
	query := `SELECT address, topic0, topic1, topic2, topic3, data, block_number, tx_hash,
			  tx_index, block_hash, log_index, removed 
			  FROM event_logs WHERE tx_hash = ? ORDER BY log_index`

	rows, err := s.db.QueryContext(ctx, query, txHash.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []*types.Log
	for rows.Next() {
		var log types.Log
		var addrStr, blockHashStr string
		var topic0Null, topic1Null, topic2Null, topic3Null, dataNull sql.NullString

		err := rows.Scan(&addrStr, &topic0Null, &topic1Null, &topic2Null, &topic3Null,
			&dataNull, &log.BlockNumber, &log.TransactionHash,
			&log.TransactionIndex, &blockHashStr, &log.LogIndex, &log.Removed)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}

		log.Address = common.HexToAddress(addrStr)
		log.BlockHash = common.HexToHash(blockHashStr)

		// Parse topics
		var topics []common.Hash
		if topic0Null.Valid {
			topics = append(topics, common.HexToHash(topic0Null.String))
		}
		if topic1Null.Valid {
			topics = append(topics, common.HexToHash(topic1Null.String))
		}
		if topic2Null.Valid {
			topics = append(topics, common.HexToHash(topic2Null.String))
		}
		if topic3Null.Valid {
			topics = append(topics, common.HexToHash(topic3Null.String))
		}
		log.Topics = topics

		if dataNull.Valid {
			log.Data = common.FromHex(dataNull.String)
		}

		logs = append(logs, &log)
	}

	return logs, nil
}
