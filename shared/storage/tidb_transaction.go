package storage

import (
	"database/sql"
	"fmt"

	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// tidbTransaction implements the Transaction interface using TiDB
type tidbTransaction struct {
	tx      *sql.Tx
	storage *TiDBStorage
}

// Commit commits the transaction
func (t *tidbTransaction) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *tidbTransaction) Rollback() error {
	return t.tx.Rollback()
}

// SetAccount sets account data in the transaction
func (t *tidbTransaction) SetAccount(address common.Address, account *types.Account) error {
	query := `INSERT INTO accounts (address, nonce, balance, code_hash, storage_root, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, NOW(), NOW())
			  ON DUPLICATE KEY UPDATE 
			  nonce = VALUES(nonce), 
			  balance = VALUES(balance), 
			  code_hash = VALUES(code_hash), 
			  storage_root = VALUES(storage_root), 
			  updated_at = NOW()`

	_, err := t.tx.Exec(query,
		address.Hex(),
		account.Nonce,
		account.Balance.String(),
		account.CodeHash.Hex(),
		account.StorageRoot.Hex())

	if err != nil {
		return fmt.Errorf("failed to set account in transaction: %w", err)
	}

	return nil
}

// SetStorageAt sets storage data in the transaction
func (t *tidbTransaction) SetStorageAt(address common.Address, key, value common.Hash) error {
	query := `INSERT INTO contract_storage (address, key, value, created_at, updated_at) 
			  VALUES (?, ?, ?, NOW(), NOW())
			  ON DUPLICATE KEY UPDATE 
			  value = VALUES(value), 
			  updated_at = NOW()`

	_, err := t.tx.Exec(query, address.Hex(), key.Hex(), value.Hex())
	if err != nil {
		return fmt.Errorf("failed to set storage in transaction: %w", err)
	}

	return nil
}

// SetCode sets contract code in the transaction
func (t *tidbTransaction) SetCode(address common.Address, code []byte) error {
	query := `INSERT INTO contract_code (address, code, code_hash, created_at, updated_at) 
			  VALUES (?, ?, ?, NOW(), NOW())
			  ON DUPLICATE KEY UPDATE 
			  code = VALUES(code), 
			  code_hash = VALUES(code_hash), 
			  updated_at = NOW()`

	codeHash := crypto.Keccak256Hash(code)
	_, err := t.tx.Exec(query, address.Hex(), common.Bytes2Hex(code), codeHash.Hex())
	if err != nil {
		return fmt.Errorf("failed to set code in transaction: %w", err)
	}

	return nil
}

// SetBlock sets block data in the transaction
func (t *tidbTransaction) SetBlock(block *types.Block) error {
	blockQuery := `INSERT INTO blocks (block_number, block_hash, parent_hash, timestamp, gas_limit, gas_used, 
						receipt_hash, state_root, created_at) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())`

	_, err := t.tx.Exec(blockQuery, block.Number, block.Hash.Hex(), block.ParentHash.Hex(),
		block.Timestamp, block.GasLimit, block.GasUsed, block.ReceiptHash.Hex(),
		block.StateRoot.Hex())
	if err != nil {
		return fmt.Errorf("failed to insert block in transaction: %w", err)
	}

	// Insert transactions
	for i, tx := range block.Transactions {
		err = t.insertTransactionInTx(tx, block.Number, uint64(i))
		if err != nil {
			return fmt.Errorf("failed to insert transaction %s in transaction: %w", tx.Hash.Hex(), err)
		}
	}

	return nil
}

// SetReceipts sets transaction receipts in the transaction
func (t *tidbTransaction) SetReceipts(receipts []*types.Receipt) error {
	for _, receipt := range receipts {
		err := t.insertReceiptInTx(receipt)
		if err != nil {
			return fmt.Errorf("failed to insert receipt %s in transaction: %w", receipt.TransactionHash.Hex(), err)
		}
	}
	return nil
}

// SetStateRoot sets state root in the transaction
func (t *tidbTransaction) SetStateRoot(blockNumber uint64, stateRoot common.Hash) error {
	query := `UPDATE blocks SET state_root = ?, updated_at = NOW() WHERE block_number = ?`

	result, err := t.tx.Exec(query, stateRoot.Hex(), blockNumber)
	if err != nil {
		return fmt.Errorf("failed to set state root in transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected in transaction: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("block not found in transaction: %d", blockNumber)
	}

	return nil
}

// insertTransactionInTx inserts a transaction within the database transaction
func (t *tidbTransaction) insertTransactionInTx(tx *types.Transaction, blockNumber uint64, txIndex uint64) error {
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

	_, err := t.tx.Exec(query, tx.Hash.Hex(), blockNumber, txIndex,
		tx.From.Hex(), toAddr, tx.Value.String(), tx.Gas,
		tx.GasPrice.String(), common.Bytes2Hex(tx.Data), tx.Nonce,
		maxFee, maxPriority, tx.V.String(), tx.R.String(), tx.S.String())

	return err
}

// insertReceiptInTx inserts a receipt within the database transaction
func (t *tidbTransaction) insertReceiptInTx(receipt *types.Receipt) error {
	query := `INSERT INTO transaction_receipts (tx_hash, tx_index, block_hash, block_number, 
			  from_address, to_address, gas_used, cumulative_gas_used, contract_address, status, 
			  effective_gas_price, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`

	var toAddr, contractAddr sql.NullString
	if receipt.To != nil {
		toAddr.String = receipt.To.Hex()
		toAddr.Valid = true
	}
	if receipt.ContractAddress != nil {
		contractAddr.String = receipt.ContractAddress.Hex()
		contractAddr.Valid = true
	}

	_, err := t.tx.Exec(query, receipt.TransactionHash.Hex(), receipt.TransactionIndex,
		receipt.BlockHash.Hex(), receipt.BlockNumber, receipt.From.Hex(),
		toAddr, receipt.GasUsed, receipt.CumulativeGasUsed, contractAddr,
		receipt.Status, receipt.EffectiveGasPrice.String())
	if err != nil {
		return fmt.Errorf("failed to insert receipt: %w", err)
	}

	// Insert logs
	for _, log := range receipt.Logs {
		err := t.insertLogInTx(log)
		if err != nil {
			return fmt.Errorf("failed to insert log: %w", err)
		}
	}

	return nil
}

// insertLogInTx inserts a log within the database transaction
func (t *tidbTransaction) insertLogInTx(log *types.Log) error {
	query := `INSERT INTO event_logs (address, topic0, topic1, topic2, topic3, data, block_number, 
			  tx_hash, tx_index, block_hash, log_index, removed, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`

	var topic0, topic1, topic2, topic3, data sql.NullString

	if len(log.Topics) > 0 {
		topic0.String = log.Topics[0].Hex()
		topic0.Valid = true
	}
	if len(log.Topics) > 1 {
		topic1.String = log.Topics[1].Hex()
		topic1.Valid = true
	}
	if len(log.Topics) > 2 {
		topic2.String = log.Topics[2].Hex()
		topic2.Valid = true
	}
	if len(log.Topics) > 3 {
		topic3.String = log.Topics[3].Hex()
		topic3.Valid = true
	}

	if len(log.Data) > 0 {
		data.String = common.Bytes2Hex(log.Data)
		data.Valid = true
	}

	_, err := t.tx.Exec(query, log.Address.Hex(), topic0, topic1, topic2, topic3, data,
		log.BlockNumber, log.TransactionHash.Hex(), log.TransactionIndex,
		log.BlockHash.Hex(), log.LogIndex, log.Removed)

	return err
}
