package test

import (
	"banking_application/api/domain/dtos"
	"banking_application/api/http/models"
	"banking_application/api/services"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

//
//type TransactionServiceTestSuite struct {
//	suite.Suite
//	mockDB        *MockDBContext
//	service       *services.TransactionService
//	mockTx        *MockTx
//	mockQueryRow  *MockRow
//	mockClient    *MockClient
//	mockValidator *MockValidator
//}
//
//// sets up the test suite
//func (suite *TransactionServiceTestSuite) SetupTest() {
//	suite.mockDB = new(MockDBContext)
//	suite.mockTx = new(MockTx)
//	suite.mockClient = new(MockClient)
//	suite.mockValidator = new(MockValidator)
//	suite.service = services.NewTransactionService(suite.mockDB, suite.mockClient, suite.mockValidator)
//
//	// Mock Begin method
//	suite.mockDB.On("Begin").Return(suite.mockTx, nil)
//	suite.mockTx.On("Rollback").Return(nil)
//}
//
//func (suite *TransactionServiceTestSuite) TestDeposit() {
//	request := dtos.DepositRequestDto{
//		AccountNumber:        "1234567890",
//		Amount:               100.0,
//		TransactionReference: "trx123",
//	}
//	var checksum = ""
//
//	account := dtos.AccountDto{
//		ID:             uuid.New().String(),
//		UserId:         uuid.New().String(),
//		AccountNumber:  request.AccountNumber,
//		AccountBalance: 500.0,
//		CheckSum:       &checksum,
//		IsActive:       true,
//	}
//
//	//checkSum, _ := suite.mockValidator.ComputeChecksum(account)
//	// Set up the mockValidator to return a valid checksum
//	suite.mockValidator.On("ComputeChecksum", account).Return("validChecksum", nil)
//	//account.CheckSum = &checkSum
//	// Set the expected checksum
//	*account.CheckSum = "validChecksum"
//
//	mockRow := new(MockRow)
//	mockRow.On("Scan", mock.AnythingOfType("*uuid.UUID"), mock.AnythingOfType("*uuid.UUID"),
//		mock.AnythingOfType("*string"), mock.AnythingOfType("*float32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*bool")).
//		Run(func(args mock.Arguments) {
//			*(args[0].(*string)) = account.ID
//			*(args[1].(*string)) = account.UserId
//			*(args[2].(*string)) = account.AccountNumber
//			*(args[3].(*float32)) = account.AccountBalance
//			*(args[4].(*string)) = *account.CheckSum
//			*(args[5].(*bool)) = account.IsActive
//		}).Return(nil)
//
//	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once()
//	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
//	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
//	suite.mockDB.On("Commit").Return(nil).Once()
//
//	// Construct a valid deposit request
//	validRequest := dtos.DepositRequestDto{
//		TransactionReference: "123456",
//		AccountNumber:        "987654321",
//		Amount:               100.0,
//	}
//
//	// Call the method being tested
//	successfulDeposit, err := suite.service.Deposit(validRequest)
//
//	// Assert expectations
//	assert.NoError(suite.T(), err)
//	assert.True(suite.T(), successfulDeposit)
//
//	// Test case: validation error
//	invalidRequest := dtos.DepositRequestDto{} // Empty request causing validation error
//	_, validationErr := suite.service.Deposit(invalidRequest)
//	assert.Error(suite.T(), validationErr)
//
//	// Test case: failure to begin transaction
//	suite.mockDB.On("Begin").Return(nil, fmt.Errorf("failed to begin transaction")).Once() // Mocking a failed transaction begin
//	_, beginErr := suite.service.Deposit(validRequest)
//	assert.Error(suite.T(), beginErr)
//
//	// Test case: duplicate transaction reference
//	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                        // Mocking a successful transaction begin
//	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once() // Mocking a successful account retrieval
//	suite.mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil).Once()    // Mocking transaction reference check
//	_, duplicateRefErr := suite.service.Deposit(validRequest)
//	assert.Error(suite.T(), duplicateRefErr)
//
//	// Test case: account not found
//	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                                   // Mocking a successful transaction begin
//	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows).Once() // Mocking account not found
//	_, accountNotFoundErr := suite.service.Deposit(validRequest)
//	assert.Error(suite.T(), accountNotFoundErr)
//
//	// Test case: checksum validation failure
//	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                        // Mocking a successful transaction begin
//	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once() // Mocking a successful account retrieval
//	suite.mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil).Once()    // Mocking checksum validation failure
//	_, checksumErr := suite.service.Deposit(validRequest)
//	assert.Error(suite.T(), checksumErr)
//
//	// Test case: failure to update account
//	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                                                                                      // Mocking a successful transaction begin
//	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once()                                                               // Mocking a successful account retrieval
//	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to update account")).Once() // Mocking a failure to update account
//	_, updateErr := suite.service.Deposit(validRequest)
//	assert.Error(suite.T(), updateErr)
//
//	// Test case: failure to insert transaction entry
//	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                                                                                                // Mocking a successful transaction begin
//	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once()                                                                         // Mocking a successful account retrieval
//	suite.mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil).Once()                                                                            // Mocking a successful account update
//	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to insert transaction entry")).Once() // Mocking a failure to insert transaction entry
//	_, insertErr := suite.service.Deposit(validRequest)
//	assert.Error(suite.T(), insertErr)
//
//	// Test case: failure to commit transaction
//	suite.mockDB.On("Begin").Return(suite.mockTx, nil).Once()                           // Mocking a successful transaction begin
//	suite.mockDB.On("QueryRow", mock.Anything, mock.Anything).Return(mockRow).Once()    // Mocking a successful account retrieval
//	suite.mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil).Once()       // Mocking a successful account update
//	suite.mockDB.On("Exec", mock.Anything, mock.Anything).Return(nil, nil).Once()       // Mocking a successful transaction entry insert
//	suite.mockDB.On("Commit").Return(fmt.Errorf("failed to commit transaction")).Once() // Mocking a failure to commit transaction
//	_, commitErr := suite.service.Deposit(validRequest)
//	assert.Error(suite.T(), commitErr)
//
//	// Assert that all expectations were met
//	suite.mockDB.AssertExpectations(suite.T())
//}
//
//// Entry point for the test suite
//func TestTransactionServiceTestSuite(t *testing.T) {
//	suite.Run(t, new(TransactionServiceTestSuite))
//}

// Mock implementations
//type MockDBContext struct {
//	mock.Mock
//}

//func (m *MockDBContext) Begin() (*sql.Tx, error) {
//	args := m.Called()
//	return args.Get(0).(*sql.Tx), args.Error(1)
//}
//
//func (m *MockDBContext) Commit() error {
//	args := m.Called()
//	return args.Error(0)
//}
//
//func (m *MockDBContext) Rollback() error {
//	args := m.Called()
//	return args.Error(0)
//}

//func (m *MockDBContext) Insert(query string, args ...interface{}) (sql.Result, error) {
//	arguments := m.Called(query, args)
//	return arguments.Get(0).(sql.Result), arguments.Error(1)
//}
//
//func (m *MockDBContext) Update(query string, args ...interface{}) (sql.Result, error) {
//	arguments := m.Called(query, args)
//	return arguments.Get(0).(sql.Result), arguments.Error(1)
//}
//
//func (m *MockDBContext) Query(query string, args ...interface{}) (*sql.Rows, error) {
//	arguments := m.Called(query, args)
//	return arguments.Get(0).(*sql.Rows), arguments.Error(1)
//}
//
//func (m *MockDBContext) QueryRow(query string, args ...interface{}) *sql.Row {
//	arguments := m.Called(query, args)
//	return arguments.Get(0).(*sql.Row)
//}

type MockClientAction struct {
	mock.Mock
}

func (m *MockClientAction) PostPayment(url string, request models.PaymentRequest) (*models.PaymentResponse, error) {
	args := m.Called(url, request)
	return args.Get(0).(*models.PaymentResponse), args.Error(1)
}

func (m *MockClientAction) GetPayment(url string) (*models.PaymentResponse, error) {
	args := m.Called(url)
	return args.Get(0).(*models.PaymentResponse), args.Error(1)
}

//type MockValidator struct {
//	mock.Mock
//}

//func (m *MockValidator) IsTransactionReferenceExist(transactionReference string) (bool, error) {
//	args := m.Called(transactionReference)
//	return args.Bool(0), args.Error(1)
//}

//func (m *MockValidator) ValidateChecksum(account dtos.AccountDto) (bool, error) {
//	args := m.Called(account)
//	return args.Bool(0), args.Error(1)
//}

//func (m *MockValidator) ComputeChecksum(account dtos.AccountDto) (string, error) {
//	args := m.Called(account)
//	return args.String(0), args.Error(1)
//}

// Test suite
type TransactionServiceTestSuite struct {
	suite.Suite
	service       *services.TransactionService
	mockDB        *MockDBContext
	mockClient    *MockClientAction
	mockValidator *MockValidator
	mockTx        *sql.Tx
}

func (suite *TransactionServiceTestSuite) SetupTest() {
	suite.mockDB = new(MockDBContext)
	suite.mockClient = new(MockClientAction)
	suite.mockValidator = new(MockValidator)
	suite.service = &services.TransactionService{
		DBContext:    suite.mockDB,
		ClientAction: suite.mockClient,
		Validator:    suite.mockValidator,
	}
	suite.mockTx = new(sql.Tx)
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

	// Validation succeeds
	suite.mockValidator.On("IsTransactionReferenceExist", request.TransactionReference).Return(false, nil)
	suite.mockValidator.On("ValidateChecksum", account).Return(true, nil)
	suite.mockValidator.On("ComputeChecksum", account).Return("validChecksum", nil)

	// Account retrieval succeeds
	mockRow := new(MockRow)
	mockRow.On("Scan", mock.AnythingOfType("*string"), mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"), mock.AnythingOfType("*float32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*bool")).
		Run(func(args mock.Arguments) {
			*(args[0].(*string)) = account.ID
			*(args[1].(*string)) = account.UserId
			*(args[2].(*string)) = account.AccountNumber
			*(args[3].(*float32)) = account.AccountBalance
			*(args[4].(*string)) = *account.CheckSum
			*(args[5].(*bool)) = account.IsActive
		}).Return(nil)

	suite.mockDB.On("QueryRow", `SELECT id, user_id, account_number, account_balance, checksum, is_active FROM accounts 
                            				WHERE account_number = $1 FOR UPDATE`, request.AccountNumber).Return(mockRow)

	// Begin transaction
	suite.mockDB.On("Begin").Return(suite.mockTx, nil)

	// Update account balance
	suite.mockDB.On("Exec", `UPDATE accounts SET account_balance = $1, updated_at = $2, checksum = $3 WHERE id = $4`,
		account.AccountBalance+request.Amount, mock.Anything, "validChecksum", account.ID).Return(nil, nil)

	// Mock third-party payment
	paymentResponse := &models.PaymentResponse{Reference: "payment123"}
	suite.mockClient.On("PostPayment", "", mock.AnythingOfType("models.PaymentRequest")).Return(paymentResponse, nil)

	// Insert transaction entry
	suite.mockDB.On("Exec", `INSERT INTO transactions (account_id, transaction_reference, transaction_type, transaction_record_type, 
                          transaction_amount, transaction_status, created_at, third_party_reference) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, account.ID, request.TransactionReference, "deposit", "credit", request.Amount, "successful", mock.Anything, "payment123").Return(nil, nil)

	// Commit transaction
	suite.mockDB.On("Commit").Return(nil)

	// Call the method being tested
	successfulDeposit, err := suite.service.Deposit(request)

	// Assert expectations
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), successfulDeposit)

	// Assert all mocks
	suite.mockDB.AssertExpectations(suite.T())
	suite.mockClient.AssertExpectations(suite.T())
	suite.mockValidator.AssertExpectations(suite.T())
}

// Add more test cases as needed...

func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}

//func (r *MockRow) Scan(dest ...interface{}) error {
//	args := r.Called(dest)
//	return args.Error(0)
//}
