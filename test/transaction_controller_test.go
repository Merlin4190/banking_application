package test

import (
	"banking_application/api/controllers"
	"banking_application/api/database"
	"banking_application/api/domain/dtos"
	htp "banking_application/api/http"
	"banking_application/api/services"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Connect to the test database
	db := database.ConnectDB()
	//defer db.Close()

	// Initialize necessary services
	client := htp.NewClient()
	validator := services.NewTransactionValidator(database.NewDBContext(db))
	transactionService := services.NewTransactionService(database.NewDBContext(db), client, validator)
	transactionController := controllers.NewTransactionController(transactionService)

	// Set up routes
	router.POST("/transaction/deposit", transactionController.Deposit())
	router.POST("/transaction/withdraw", transactionController.Withdraw())

	return router
}

func TestDepositRequest(t *testing.T) {
	router := setupRouter()
	depositRequest := dtos.DepositRequestDto{
		Amount:               100.0,
		AccountNumber:        "079789496",
		TransactionReference: "123456789123456",
	}

	t.Run("negative deposit amount test", func(t *testing.T) {
		requestBody, _ := json.Marshal(depositRequest)
		req, _ := http.NewRequest("POST", "/transaction/deposit", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectedResponse := `{"error":"amount cannot be less than or equal zero","success":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})

	t.Run("duplicate deposit test", func(t *testing.T) {
		requestBody, _ := json.Marshal(depositRequest)
		req, _ := http.NewRequest("POST", "/transaction/deposit", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectedResponse := `{"error":"duplicate transaction","success":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})

	t.Run("bad request with invalid JSON", func(t *testing.T) {
		invalidRequestBody := `{"amount":"invalid_amount","account_number":"1234567890","transaction_reference":"TX12345"}`
		req, _ := http.NewRequest("POST", "/transaction/deposit", bytes.NewBufferString(invalidRequestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"json: cannot unmarshal`)
		assert.Contains(t, w.Body.String(), `"success":false`)
	})

	t.Run("deposit service error", func(t *testing.T) {
		requestBody, _ := json.Marshal(depositRequest)
		req, _ := http.NewRequest("POST", "/transaction/deposit", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expectedResponse := `{"error":"internal server error","success":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})
}

func TestWithdrawRequest(t *testing.T) {
	router := setupRouter()
	withdrawRequest := dtos.WithdrawRequestDto{
		Amount:               -50000000.0,
		AccountNumber:        "0797894906",
		TransactionReference: "123456789471379",
	}

	t.Run("successful withdrawal", func(t *testing.T) {
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequest("POST", "/transaction/withdraw", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expectedResponse := `{"message":"withdrawal transaction successful","success":true}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})

	t.Run("bad request with invalid JSON", func(t *testing.T) {
		invalidRequestBody := `{"amount":"invalid_amount","account_number":"1234567890","transaction_reference":"TX67890"}`
		req, _ := http.NewRequest("POST", "/transaction/withdraw", bytes.NewBufferString(invalidRequestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"json: cannot unmarshal`)
		assert.Contains(t, w.Body.String(), `"success":false`)
	})

	t.Run("negative withdrawal amount test", func(t *testing.T) {
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequest("POST", "/transaction/withdraw", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectedResponse := `{"error":"amount cannot be less than or equal zero","success":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})

	t.Run("insufficient balance test", func(t *testing.T) {
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequest("POST", "/transaction/withdraw", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expectedResponse := `{"error":"insufficient balance for withdrawal","success":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})

	t.Run("duplicate deposit test", func(t *testing.T) {
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequest("POST", "/transaction/withdraw", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		expectedResponse := `{"error":"duplicate transaction","success":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})

	t.Run("withdrawal service error", func(t *testing.T) {
		requestBody, _ := json.Marshal(withdrawRequest)
		req, _ := http.NewRequest("POST", "/transaction/withdraw", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		expectedResponse := `{"error":"internal server error","success":false}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})
}
