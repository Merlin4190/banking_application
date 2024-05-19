package test

import (
	"banking_application/api/domain/dtos"
	"banking_application/api/services"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TransactionServiceTestSuite struct {
	suite.Suite
	mockDB        *MockDBContext
	service       *services.TransactionService
	mockTx        *MockTx
	mockQueryRow  *MockRow
	mockClient    *MockClient
	mockValidator *MockValidator
}

// sets up the test suite
func (suite *TransactionServiceTestSuite) SetupTest() {
	suite.mockDB = new(MockDBContext)
	suite.mockTx = new(MockTx)
	suite.mockClient = new(MockClient)
	suite.mockValidator = new(MockValidator)
	suite.service = services.NewTransactionService(suite.mockDB, suite.mockClient, suite.mockValidator)

	// Mock Begin method
	suite.mockDB.On("Begin").Return(suite.mockTx, nil)
	suite.mockTx.On("Rollback").Return(nil)
}

func (suite *TransactionServiceTestSuite) TestDeposit() {
	request := dtos.DepositRequestDto{
		AccountNumber:        "1234567890",
		Amount:               100.0,
		TransactionReference: "trx123",
	}
	var checksum = ""

	account := dtos.AccountDto{
		ID:             uuid.New().String(),
		UserId:         uuid.New().String(),
		AccountNumber:  request.AccountNumber,
		AccountBalance: 500.0,
		CheckSum:       &checksum,
		IsActive:       true,
	}

	checkSum, _ := suite.mockValidator.ComputeChecksum(account)

	account.CheckSum = &checkSum

	mockRow := new(MockRow)
	mockRow.On("Scan", mock.AnythingOfType("*uuid.UUID"), mock.AnythingOfType("*uuid.UUID"),
		mock.AnythingOfType("*string"), mock.AnythingOfType("*float32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*bool")).
		Run(func(args mock.Arguments) {
			*(args[0].(*string)) = account.ID
			*(args[1].(*string)) = account.UserId
			*(args[2].(*string)) = account.AccountNumber
			*(args[3].(*float32)) = account.AccountBalance
			*(args[4].(*string)) = *account.CheckSum
			*(args[5].(*bool)) = account.IsActive
		}).Return(nil)

	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once()
	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
	suite.mockDB.On("Commit").Return(nil).Once()

	// Construct a valid deposit request
	validRequest := dtos.DepositRequestDto{
		TransactionReference: "123456",
		AccountNumber:        "987654321",
		Amount:               100.0,
	}

	// Call the method being tested
	successfulDeposit, err := suite.service.Deposit(validRequest)

	// Assert expectations
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), successfulDeposit)

	// Test case: validation error
	invalidRequest := dtos.DepositRequestDto{} // Empty request causing validation error
	_, validationErr := suite.service.Deposit(invalidRequest)
	assert.Error(suite.T(), validationErr)

	// Test case: failure to begin transaction
	suite.mockDB.On("Begin").Return(nil, fmt.Errorf("failed to begin transaction")).Once() // Mocking a failed transaction begin
	_, beginErr := suite.service.Deposit(validRequest)
	assert.Error(suite.T(), beginErr)

	// Test case: duplicate transaction reference
	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                        // Mocking a successful transaction begin
	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once() // Mocking a successful account retrieval
	suite.mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil).Once()    // Mocking transaction reference check
	_, duplicateRefErr := suite.service.Deposit(validRequest)
	assert.Error(suite.T(), duplicateRefErr)

	// Test case: account not found
	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                                   // Mocking a successful transaction begin
	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows).Once() // Mocking account not found
	_, accountNotFoundErr := suite.service.Deposit(validRequest)
	assert.Error(suite.T(), accountNotFoundErr)

	// Test case: checksum validation failure
	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                        // Mocking a successful transaction begin
	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once() // Mocking a successful account retrieval
	suite.mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil).Once()    // Mocking checksum validation failure
	_, checksumErr := suite.service.Deposit(validRequest)
	assert.Error(suite.T(), checksumErr)

	// Test case: failure to update account
	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                                                                                      // Mocking a successful transaction begin
	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once()                                                               // Mocking a successful account retrieval
	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to update account")).Once() // Mocking a failure to update account
	_, updateErr := suite.service.Deposit(validRequest)
	assert.Error(suite.T(), updateErr)

	// Test case: failure to insert transaction entry
	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                                                                                                // Mocking a successful transaction begin
	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once()                                                                         // Mocking a successful account retrieval
	suite.mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil).Once()                                                                            // Mocking a successful account update
	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to insert transaction entry")).Once() // Mocking a failure to insert transaction entry
	_, insertErr := suite.service.Deposit(validRequest)
	assert.Error(suite.T(), insertErr)

	// Test case: failure to commit transaction
	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                           // Mocking a successful transaction begin
	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once()    // Mocking a successful account retrieval
	suite.mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil).Once()       // Mocking a successful account update
	suite.mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil).Once()       // Mocking a successful transaction entry insert
	suite.mockDB.On("Commit").Return(fmt.Errorf("failed to commit transaction")).Once() // Mocking a failure to commit transaction
	_, commitErr := suite.service.Deposit(validRequest)
	assert.Error(suite.T(), commitErr)

	// Assert that all expectations were met
	suite.mockDB.AssertExpectations(suite.T())
}

// Entry point for the test suite
func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}
