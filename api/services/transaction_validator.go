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
)

type Validator interface {
	IsTransactionReferenceExist(transactionReference string) (bool, error)
	ValidateChecksum(account dtos.AccountDto) (bool, error)
	ComputeChecksum(account dtos.AccountDto) (string, error)
}

type TransactionValidator struct {
	Validator *TransactionValidator
	ctx       database.DBContext
}

func NewTransactionValidator(ctx database.DBContext) *TransactionValidator {
	return &TransactionValidator{ctx: ctx}
}

func (v *TransactionValidator) IsTransactionReferenceExist(transactionReference string) (bool, error) {
	return isTransactionReferenceExist(v.ctx, transactionReference)
}

func (v *TransactionValidator) ValidateChecksum(account dtos.AccountDto) (bool, error) {
	return validateChecksum(account)
}

func (v *TransactionValidator) ComputeChecksum(account dtos.AccountDto) (string, error) {
	return computeChecksum(account)
}

func isTransactionReferenceExist(ctx database.DBContext, transactionRef string) (bool, error) {
	// Check if the transaction reference already exists in the database
	var id string
	err := ctx.QueryRow(`SELECT id FROM transactions WHERE transaction_reference = $1`, transactionRef).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func computeChecksum(account dtos.AccountDto) (string, error) {
	// Concatenate the data
	data := fmt.Sprintf("%s|%s|%.4f|%t|%s", account.UserId, account.AccountNumber, account.AccountBalance, account.IsActive, account.ID)

	// Calculate the checksum
	key := "9b4e7a91c3b9e5a9a5b3c2e80e5d2a%v"
	checkSumValue, err := util.AESEncrypt(data, key)
	if err != nil {
		return "", fmt.Errorf("error computing checksum: %v", err)
	}
	return checkSumValue, nil
}

func validateChecksum(account dtos.AccountDto) (bool, error) {
	isValid := false

	// Attempt to decrypt the checksum
	key := "9b4e7a91c3b9e5a9a5b3c2e80e5d2a%v"
	decryptedValue, err := util.AESDecrypt(*account.CheckSum, key)
	if err != nil {
		return false, err
	}

	// Split the decrypted value into parts
	values := strings.Split(decryptedValue, "|")
	if len(values) == 5 {
		// Extract values from the decrypted data
		accountHolderID := values[0]
		accountNumber := values[1]
		balance, _ := strconv.ParseFloat(values[2], 32)
		status := values[3]
		//accountID := values[4]

		accountBalance := float32(balance)
		/*userId, _ := uuid.Parse(accountHolderID)
		accountUUID, _ := uuid.Parse(accountID)*/

		// Validate against wallet attributes
		if /*account.ID == accountID &&*/
		account.UserId == accountHolderID &&
			accountNumber == account.AccountNumber &&
			accountBalance == account.AccountBalance &&
			status == strconv.FormatBool(account.IsActive) {
			isValid = true
		}
	}

	return isValid, nil
}
