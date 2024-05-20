package test

import (
	"banking_application/api/controllers"
	"banking_application/api/domain/dtos"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockTransactionService is a mock implementation of the Transaction interface
type MockTransactionService struct {
	forceError bool
}

func (m *MockTransactionService) Deposit(account dtos.DepositRequestDto) (bool, error) {
	if m.forceError {
		return false, errors.New("internal server error")
	}
	return true, nil
}

func (m *MockTransactionService) ForceError(force bool) {
	m.forceError = force
}

func (m *MockTransactionService) Withdraw(request dtos.WithdrawRequestDto) (bool, error) {
	return false, nil // Implement as needed for other tests
}

func (m *MockTransactionService) Transfer(request dtos.TransferRequestDto) (bool, error) {
	return false, nil // Implement as needed for other tests
}

func TestDeposit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockService := &MockTransactionService{}
	transactionController := controllers.NewTransactionController(mockService)

	router.POST("/deposit", transactionController.Deposit())

	t.Run("successful deposit", func(t *testing.T) {
		depositRequest := dtos.DepositRequestDto{
			Amount:               100.0,
			AccountNumber:        "1234567890",
			TransactionReference: "TX12345",
		}
		requestBody, _ := json.Marshal(depositRequest)
		req, _ := http.NewRequest("POST", "/deposit", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expectedResponse := `{"message":"deposit transaction successful","success":true}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})

	t.Run("bad request with invalid JSON", func(t *testing.T) {
		invalidRequestBody := `{"amount":"invalid_amount","account_number":"1234567890","transaction_reference":"TX12345"}`
		req, _ := http.NewRequest("POST", "/deposit", bytes.NewBufferString(invalidRequestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"json: cannot unmarshal`)
		assert.Contains(t, w.Body.String(), `"success":false`)
	})

	t.Run("deposit service error", func(t *testing.T) {
		mockService.ForceError(true)

		depositRequest := dtos.DepositRequestDto{
			Amount:               100.0,
			AccountNumber:        "1234567890",
			TransactionReference: "TX12345",
		}
		requestBody, _ := json.Marshal(depositRequest)
		req, _ := http.NewRequest("POST", "/deposit", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expectedResponse := `{"error":"internal server error","success":false}` // Adjust based on your actual error response
		assert.JSONEq(t, expectedResponse, w.Body.String())

		mockService.ForceError(false)
	})
}

func TestWithdraw(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockService := &MockTransactionService{}
	transactionController := controllers.NewTransactionController(mockService)

	router.POST("/withdraw", transactionController.Withdraw())

	t.Run("successful withdrawal", func(t *testing.T) {
		withdrawRequest := dtos.WithdrawRequestDto{
			Amount:               50.0,
			AccountNumber:        "1234567890",
			TransactionReference: "TX67890",
		}
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequest("POST", "/withdraw", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expectedResponse := `{"message":"withdrawal transaction successful","success":true}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})

	t.Run("bad request with invalid JSON", func(t *testing.T) {
		invalidRequestBody := `{"amount":"invalid_amount","account_number":"1234567890","transaction_reference":"TX67890"}`
		req, _ := http.NewRequest("POST", "/withdraw", bytes.NewBufferString(invalidRequestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"json: cannot unmarshal`)
		assert.Contains(t, w.Body.String(), `"success":false`)
	})

	t.Run("withdrawal service error", func(t *testing.T) {
		mockService.ForceError(true)

		withdrawRequest := dtos.WithdrawRequestDto{
			Amount:               50.0,
			AccountNumber:        "1234567890",
			TransactionReference: "TX67890",
		}
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequest("POST", "/withdraw", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expectedResponse := `{"error":"internal server error","success":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())

		mockService.ForceError(false)
	})
}
