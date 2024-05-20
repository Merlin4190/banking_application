package main

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

// Mock implementations
type MockDBContext struct {
	mock.Mock
}

func (m *MockDBContext) Begin() (*sql.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockDBContext) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDBContext) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDBContext) Insert(query string, args ...interface{}) (sql.Result, error) {
	arguments := m.Called(append([]interface{}{query}, args...)...)
	return arguments.Get(0).(sql.Result), arguments.Error(1)
}

func (m *MockDBContext) Update(query string, args ...interface{}) (sql.Result, error) {
	arguments := m.Called(append([]interface{}{query}, args...)...)
	return arguments.Get(0).(sql.Result), arguments.Error(1)
}

func (m *MockDBContext) Query(query string, args ...interface{}) (*sql.Rows, error) {
	arguments := m.Called(append([]interface{}{query}, args...)...)
	return arguments.Get(0).(*sql.Rows), arguments.Error(1)
}

func (m *MockDBContext) QueryRow(query string, args ...interface{}) *sql.Row {
	arguments := m.Called(append([]interface{}{query}, args...)...)
	return arguments.Get(0).(*sql.Row)
}

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

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) IsTransactionReferenceExist(transactionReference string) (bool, error) {
	args := m.Called(transactionReference)
	return args.Bool(0), args.Error(1)
}

func (m *MockValidator) ValidateChecksum(account dtos.AccountDto) (bool, error) {
	args := m.Called(account)
	return args.Bool(0), args.Error(1)
}

func (m *MockValidator) ComputeChecksum(account dtos.AccountDto) (string, error) {
	args := m.Called(account)
	return args.String(0), args.Error(1)
}

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
	//suite.mockDB.On("QueryRow", `SELECT id, user_id, account_number, account_balance, checksum, is_active FROM accounts
	//                        				WHERE account_number = $1 FOR UPDATE`, request.AccountNumber).Return(mockRow)
	//suite.mockDB.On("QueryRow", `SELECT id, user_id, account_number, account_balance, checksum, is_active FROM accounts WHERE account_number = $1 FOR UPDATE`, request.AccountNumber).Return(&sql.Row{})
	//mockRow := &sql.Row{}
	suite.mockDB.On("QueryRow", `SELECT id, user_id, account_number, account_balance, checksum, is_active FROM accounts WHERE account_number = $1 FOR UPDATE`, request.AccountNumber).Return(mockRow)

	// Mocking the Scan call on the returned row
	mockRow.On("Scan", &account.ID, &account.UserId, &account.AccountNumber, &account.AccountBalance, &account.CheckSum, &account.IsActive).Return(nil).Run(func(args mock.Arguments) {
		*(args[0].(*string)) = account.ID
		*(args[1].(*string)) = account.UserId
		*(args[2].(*string)) = account.AccountNumber
		*(args[3].(*float32)) = account.AccountBalance
		*(args[4].(*string)) = *account.CheckSum
		*(args[5].(*bool)) = account.IsActive
	})

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

func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}

type MockRow struct {
	mock.Mock
}

func (r *MockRow) Scan(dest ...interface{}) error {
	args := r.Called(dest...)
	return args.Error(0)
}
