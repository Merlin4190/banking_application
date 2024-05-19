package services

import (
	"banking_application/api/database"
	"banking_application/api/domain/dtos"
	"banking_application/api/http"
	"banking_application/api/http/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Transaction interface {
	Deposit(request dtos.DepositRequestDto) (bool, error)
	Withdraw(request dtos.WithdrawRequestDto) (bool, error)
	Transfer(request dtos.TransferRequestDto) (bool, error)
}

type TransactionService struct {
	DBContext    database.DBContext
	ClientAction http.ClientAction
	Validator    Validator
}

func NewTransactionService(dbContext database.DBContext, clientAction http.ClientAction, validator Validator) *TransactionService {
	return &TransactionService{DBContext: dbContext, ClientAction: clientAction, Validator: validator}
}

func (s *TransactionService) Deposit(request dtos.DepositRequestDto) (bool, error) {
	validationErr := validate.Struct(request)

	if validationErr != nil {
		return false, fmt.Errorf("validation %v", validationErr)
	}

	//Check for duplicate transaction
	isExist, refErr := s.Validator.IsTransactionReferenceExist(request.TransactionReference)
	if refErr != nil {
		return false, fmt.Errorf("transaction reference validation failed: %v", refErr)
	}
	if isExist {
		return false, fmt.Errorf("duplicate transaction")
	}

	var account dtos.AccountDto

	//Validate AccountNumber
	if err := s.DBContext.QueryRow(`SELECT id, user_id, account_number, account_balance, checksum, is_active FROM accounts 
                            				WHERE account_number = $1 FOR UPDATE`, request.AccountNumber).Scan(&account.ID, &account.UserId,
		&account.AccountNumber, &account.AccountBalance, &account.CheckSum, &account.IsActive); err != nil {
		// Check if no rows were returned
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("account not found for accountnumber %s", request.AccountNumber)
		}
		return false, err
	}

	//Validate Checksum
	isValid, checksumErr := s.Validator.ValidateChecksum(account)
	if checksumErr != nil {
		return false, fmt.Errorf("checksum validation failed: %v", checksumErr)
	}
	if !isValid {
		return false, fmt.Errorf("account %s is locked and cannot carry out this transaction", request.AccountNumber)
	}

	// Begin transaction
	tx, err := s.DBContext.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	account.AccountBalance += request.Amount
	account.UpdatedAt = time.Now()
	//Compute Checksum
	checksum, _ := s.Validator.ComputeChecksum(account)

	account.CheckSum = &checksum

	//Update Account
	_, updateErr := tx.Exec(`UPDATE accounts SET account_balance = $1, updated_at = $2, checksum = $3 WHERE id = $4`,
		account.AccountBalance, account.UpdatedAt, account.CheckSum, account.ID)
	if updateErr != nil {
		return false, fmt.Errorf("transaction failed: %v", updateErr)
	}

	//Call the third-party api for payment
	paymentRequest := models.PaymentRequest{
		AccountID: account.AccountNumber,
		Reference: request.TransactionReference,
		Amount:    request.Amount,
	}

	response, clientErr := s.ClientAction.PostPayment("", paymentRequest)
	if clientErr != nil {
		return false, fmt.Errorf("payment failed: %v", clientErr)
	}

	//Create Transaction entry
	_, insertErr := tx.Exec(`INSERT INTO transactions (account_id, transaction_reference, transaction_type, transaction_record_type, 
                          transaction_amount, transaction_status, created_at, third_party_reference) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, account.ID, request.TransactionReference, "deposit", "credit", request.Amount, "successful", time.Now(), response.Reference)

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
	validationErr := validate.Struct(request)
	if validationErr != nil {
		return false, fmt.Errorf("validation error: %v", validationErr)
	}

	//Check for duplicate transaction
	isExist, refErr := s.Validator.IsTransactionReferenceExist(request.TransactionReference)
	if refErr != nil {
		return false, fmt.Errorf("transaction reference validation failed: %v", refErr)
	}
	if isExist {
		return false, fmt.Errorf("duplicate transaction")
	}

	var account dtos.AccountDto

	// Validate AccountNumber and lock the row
	if err := s.DBContext.QueryRow(`
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
	isValid, checksumErr := s.Validator.ValidateChecksum(account)
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

	// Begin transaction
	tx, err := s.DBContext.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Update Account balance
	account.AccountBalance -= request.Amount
	account.UpdatedAt = time.Now()

	// Compute new checksum
	checksum, checksumErr := s.Validator.ComputeChecksum(account)
	if checksumErr != nil {
		return false, fmt.Errorf("checksum computation failed: %v", checksumErr)
	}
	account.CheckSum = &checksum

	// Update Account in the database
	_, updateErr := tx.Exec(`
		UPDATE accounts 
		SET account_balance = $1, updated_at = $2, checksum = $3
		WHERE id = $4`, account.AccountBalance, account.UpdatedAt, account.CheckSum, account.ID)
	if updateErr != nil {
		return false, fmt.Errorf("update account failed: %v", updateErr)
	}

	//Call the third-party api for payment
	paymentRequest := models.PaymentRequest{
		AccountID: account.AccountNumber,
		Reference: request.TransactionReference,
		Amount:    request.Amount,
	}

	response, clientErr := s.ClientAction.PostPayment("", paymentRequest)
	if clientErr != nil {
		return false, fmt.Errorf("payment failed: %v", clientErr)
	}

	// Create Transaction entry
	_, insertErr := tx.Exec(`
		INSERT INTO transactions (account_id, transaction_reference, transaction_type, transaction_record_type, 
			transaction_amount, transaction_status, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		account.ID, response.Reference, "withdrawal", "debit", response.Amount, "successful", time.Now())
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
