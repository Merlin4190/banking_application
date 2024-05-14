package test

import (
	"banking_application/api/database"
	"banking_application/api/domain/dtos"
	"banking_application/api/services"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRow struct {
	mock.Mock
	err error
}

func (r *MockRow) Scan(dest ...interface{}) error {
	args := r.Called(dest)
	if r.err != nil {
		return r.err
	}
	for i, d := range dest {
		*(d.(*interface{})) = args.Get(i)
	}
	return nil
}

func TestDeposit(t *testing.T) {
	// Initialize the mock DB context
	mockDB := database.MockDBContext{}

	// Initialize the TransactionService with the mock DB context
	transactionService := &services.TransactionService{DBContext: &mockDB}

	tx := new(sql.Tx)
	mockDB.On("Begin").Return(tx, nil)

	request := dtos.DepositRequestDto{
		AccountNumber:        "1234567890",
		Amount:               100.0,
		TransactionReference: "trx123",
	}

	account := dtos.AccountDto{
		ID:             uuid.New(),
		UserId:         uuid.New(),
		AccountNumber:  request.AccountNumber,
		AccountBalance: 500.0,
		CheckSum:       "some-checksum",
		IsActive:       true,
	}

	mockRow := new(MockRow)
	mockRow.On("Scan", mock.Anything).Return(func(dest ...interface{}) error {
		dest[0] = account.ID
		dest[1] = account.UserId
		dest[2] = account.AccountNumber
		dest[3] = account.AccountBalance
		dest[4] = account.CheckSum
		dest[5] = account.IsActive
		return nil
	})

	// Test case: successful deposit
	mockDB.On("Begin").Return(nil).Once()                                                                                                                                            // Mocking a successful transaction begin
	mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow, nil).Once()                                                                                                  // Mocking a successful account retrieval
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()                // Mocking a successful account update
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once() // Mocking a successful transaction entry insert
	mockDB.On("Commit").Return(nil).Once()                                                                                                                                           // Mocking a successful transaction commit

	// Construct a valid deposit request
	validRequest := dtos.DepositRequestDto{
		TransactionReference: "123456",
		AccountNumber:        "987654321",
		Amount:               100.0,
	}

	// Call the method being tested
	successfulDeposit, err := transactionService.Deposit(validRequest)

	// Assert expectations
	assert.NoError(t, err)
	assert.True(t, successfulDeposit)

	// Test case: validation error
	invalidRequest := dtos.DepositRequestDto{} // Empty request causing validation error
	_, validationErr := transactionService.Deposit(invalidRequest)
	assert.Error(t, validationErr)

	// Test case: failure to begin transaction
	mockDB.On("Begin").Return(fmt.Errorf("failed to begin transaction")).Once() // Mocking a failed transaction begin
	_, beginErr := transactionService.Deposit(validRequest)
	assert.Error(t, beginErr)

	// Test case: duplicate transaction reference
	mockDB.On("Begin").Return(nil).Once()                                                           // Mocking a successful transaction begin
	mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow, nil).Once()                 // Mocking a successful account retrieval
	mockDB.On("IsTransactionReferenceExist", mock.Anything, mock.Anything).Return(true, nil).Once() // Mocking a duplicate transaction reference
	_, duplicateRefErr := transactionService.Deposit(validRequest)
	assert.Error(t, duplicateRefErr)

	// Test case: account not found
	mockDB.On("Begin").Return(nil).Once()                                                 // Mocking a successful transaction begin
	mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows).Once() // Mocking account not found
	_, accountNotFoundErr := transactionService.Deposit(validRequest)
	assert.Error(t, accountNotFoundErr)

	// Test case: checksum validation failure
	mockDB.On("Begin").Return(nil).Once()                                           // Mocking a successful transaction begin
	mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow, nil).Once() // Mocking a successful account retrieval
	mockDB.On("ValidateChecksum", mock.Anything).Return(false, nil).Once()          // Mocking a checksum validation failure
	_, checksumErr := transactionService.Deposit(validRequest)
	assert.Error(t, checksumErr)

	// Test case: failure to update account
	mockDB.On("Begin").Return(nil).Once()                                                                                                                                                                // Mocking a successful transaction begin
	mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow, nil).Once()                                                                                                                      // Mocking a successful account retrieval
	mockDB.On("ValidateChecksum", mock.Anything).Return(true, nil).Once()                                                                                                                                // Mocking a successful checksum validation
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to update account")).Once() // Mocking a failure to update account
	_, updateErr := transactionService.Deposit(validRequest)
	assert.Error(t, updateErr)

	// Test case: failure to insert transaction entry
	mockDB.On("Begin").Return(nil).Once()                                                                                                                                                                                         // Mocking a successful transaction begin
	mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow, nil).Once()                                                                                                                                               // Mocking a successful account retrieval
	mockDB.On("ValidateChecksum", mock.Anything).Return(true, nil).Once()                                                                                                                                                         // Mocking a successful checksum validation
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()                                                             // Mocking a successful account update
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to insert transaction entry")).Once() // Mocking a failure to insert transaction entry
	_, insertErr := transactionService.Deposit(validRequest)
	assert.Error(t, insertErr)

	// Test case: failure to commit transaction
	mockDB.On("Begin").Return(nil).Once()                                                                                                                                            // Mocking a successful transaction begin
	mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow, nil).Once()                                                                                                  // Mocking a successful account retrieval
	mockDB.On("ValidateChecksum", mock.Anything).Return(true, nil).Once()                                                                                                            // Mocking a successful checksum validation
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()                // Mocking a successful account update
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once() // Mocking a successful transaction entry insert
	mockDB.On("Commit").Return(fmt.Errorf("failed to commit transaction")).Once()                                                                                                    // Mocking a failure to commit transaction
	_, commitErr := transactionService.Deposit(validRequest)
	assert.Error(t, commitErr)

	// Assert that all expectations were met
	mockDB.AssertExpectations(t)
}
