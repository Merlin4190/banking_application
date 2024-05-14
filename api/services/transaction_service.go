package services

import (
	"banking_application/api/database"
	"banking_application/api/domain/dtos"
	"banking_application/api/util"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ITransactionService interface {
	Deposit(request dtos.DepositRequestDto) (bool, error)
	Withdraw(request dtos.WithdrawRequestDto) (bool, error)
	Transfer(request dtos.TransferRequestDto) (bool, error)
}

type TransactionService struct {
	DBContext database.DBContext
}

func NewTransactionService(dbContext database.DBContext) *TransactionService {
	return &TransactionService{DBContext: dbContext}
}

func (s *TransactionService) Deposit(request dtos.DepositRequestDto) (bool, error) {
	// Begin transaction
	tx, err := s.DBContext.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	validationErr := validate.Struct(request)

	if validationErr != nil {
		return false, fmt.Errorf("validation %v", validationErr)
	}

	//Check for duplicate transaction
	isExist, refErr := IsTransactionReferenceExist(request.TransactionReference, *tx)
	if refErr != nil {
		return false, fmt.Errorf("transaction reference validation failed: %v", refErr)
	}
	if isExist {
		return false, fmt.Errorf("duplicate transaction")
	}

	var account dtos.AccountDto

	//Validate AccountNumber
	if err := tx.QueryRow(`SELECT id, user_id, account_number, account_balance, checksum, is_active FROM accounts 
                            				WHERE account_number = $1 FOR UPDATE`, request.AccountNumber).Scan(&account.ID, &account.UserId,
		&account.AccountNumber, &account.AccountBalance, &account.CheckSum, &account.IsActive); err != nil {
		// Check if no rows were returned
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("account not found for accountnumber %s", request.AccountNumber)
		}
		return false, err
	}

	//Validate Checksum
	isValid, checksumErr := ValidateChecksum(account)
	if checksumErr != nil {
		return false, fmt.Errorf("checksum validation failed: %v", checksumErr)
	}
	if !isValid {
		return false, fmt.Errorf("account %s is locked and cannot carry out this transaction", request.AccountNumber)
	}

	account.AccountBalance += request.Amount
	account.UpdatedAt = time.Now()
	//Compute Checksum
	checksum, _ := ComputeChecksum(account)

	account.CheckSum = checksum

	//Update Account
	_, updateErr := tx.Exec(`UPDATE accounts SET account_balance = $1, updated_at = $2, checksum = $3 WHERE id = $4`,
		account.AccountBalance, account.UpdatedAt, account.CheckSum, account.ID)
	if updateErr != nil {
		return false, fmt.Errorf("transaction failed: %v", updateErr)
	}

	//Create Transaction entry
	_, insertErr := tx.Exec(`INSERT INTO transactions (account_id, transaction_reference, transaction_type, transaction_record_type, 
                          transaction_amount, transaction_status, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`, account.ID, request.TransactionReference, "deposit", "credit", request.Amount, "successful", time.Now())

	if insertErr != nil {
		return false, fmt.Errorf("transaction entry failed: %v", insertErr)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return true, nil
}

func (s *TransactionService) Withdraw(request dtos.WithdrawRequestDto) (bool, error) {
	// Begin transaction
	tx, err := s.DBContext.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	validationErr := validate.Struct(request)
	if validationErr != nil {
		return false, fmt.Errorf("validation error: %v", validationErr)
	}

	//Check for duplicate transaction
	isExist, refErr := IsTransactionReferenceExist(request.TransactionReference, *tx)
	if refErr != nil {
		return false, fmt.Errorf("transaction reference validation failed: %v", refErr)
	}
	if isExist {
		return false, fmt.Errorf("duplicate transaction")
	}

	var account dtos.AccountDto

	// Validate AccountNumber and lock the row
	if err := tx.QueryRow(`
		SELECT id, user_id, account_number, account_balance, checksum, is_active
		FROM accounts 
		WHERE account_number = $1 FOR UPDATE`, request.AccountNumber).Scan(
		&account.ID, &account.UserId, &account.AccountNumber,
		&account.AccountBalance, &account.CheckSum, &account.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("account not found for accountnumber %s", request.AccountNumber)
		}
		return false, err
	}

	// Validate Checksum
	isValid, checksumErr := ValidateChecksum(account)
	if checksumErr != nil {
		return false, fmt.Errorf("checksum validation failed: %v", checksumErr)
	}
	if !isValid {
		return false, fmt.Errorf("account %s is locked and cannot carry out this transaction", request.AccountNumber)
	}

	// Check if the account has sufficient balance for withdrawal
	if account.AccountBalance < request.Amount {
		return false, fmt.Errorf("insufficient balance for withdrawal")
	}

	// Update Account balance
	account.AccountBalance -= request.Amount
	account.UpdatedAt = time.Now()

	// Compute new checksum
	checksum, checksumErr := ComputeChecksum(account)
	if checksumErr != nil {
		return false, fmt.Errorf("checksum computation failed: %v", checksumErr)
	}
	account.CheckSum = checksum

	// Update Account in the database
	_, updateErr := tx.Exec(`
		UPDATE accounts 
		SET account_balance = $1, updated_at = $2, checksum = $3
		WHERE id = $4`, account.AccountBalance, account.UpdatedAt, account.CheckSum, account.ID)
	if updateErr != nil {
		return false, fmt.Errorf("update account failed: %v", updateErr)
	}

	// Create Transaction entry
	_, insertErr := tx.Exec(`
		INSERT INTO transactions (account_id, transaction_reference, transaction_type, transaction_record_type, 
			transaction_amount, transaction_status, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		account.ID, request.TransactionReference, "withdrawal", "debit", request.Amount, "successful", time.Now())
	if insertErr != nil {
		return false, fmt.Errorf("insert transaction entry failed: %v", insertErr)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return true, nil
}

func (s *TransactionService) Transfer(request dtos.TransferRequestDto) (bool, error) {
	return false, nil
}

func IsTransactionReferenceExist(transactionRef string, tx sql.Tx) (bool, error) {
	// Check if the transaction reference already exists in the database
	var id string
	err := tx.QueryRow(`SELECT id FROM transactions WHERE transaction_reference = $1`, transactionRef).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ComputeChecksum(account dtos.AccountDto) (string, error) {
	// Concatenate the data
	data := fmt.Sprintf("%d|%s|%.4f|%t|%d", account.UserId, account.AccountNumber, account.AccountBalance, account.IsActive, account.ID)

	// Calculate the checksum (not implemented in this example)
	key := ""
	checkSumValue, err := util.AESEncrypt(data, key)
	if err != nil {
		return "", fmt.Errorf("error computing checksum: %v", err)
	}
	return checkSumValue, nil
}

func ValidateChecksum(account dtos.AccountDto) (bool, error) {
	isValid := false

	// Attempt to decrypt the checksum
	key := ""
	decryptedValue, err := util.AESDecrypt(account.CheckSum, key)
	if err != nil {
		return false, err
	}

	// Split the decrypted value into parts
	values := strings.Split(decryptedValue, "|")
	if len(values) == 5 {
		// Extract values from the decrypted data
		accountHolderID, _ := strconv.Atoi(values[0])
		accountNumber := values[1]
		balance, _ := strconv.ParseFloat(values[2], 32)
		status := values[3]
		accountID := values[4]

		accountBalance := float32(balance)

		// Validate against wallet attributes
		if fmt.Sprintf("%d", account.ID) == accountID &&
			fmt.Sprintf("%d", account.UserId) == fmt.Sprintf("%d", accountHolderID) &&
			accountNumber == account.AccountNumber &&
			accountBalance == account.AccountBalance &&
			status == strconv.FormatBool(account.IsActive) {
			isValid = true
		}
	}

	return isValid, nil
}
