package controllers

import (
	"banking_application/api/domain/dtos"
	"banking_application/api/services"
	"banking_application/api/util"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type TransactionController struct {
	transactionService services.Transaction
}

func NewTransactionController(service services.Transaction) *TransactionController {
	return &TransactionController{transactionService: service}
}

func (s *TransactionController) Deposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		account := dtos.DepositRequestDto{}

		if err := c.BindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
			return
		}

		_, serviceErr := s.transactionService.Deposit(account)

		if serviceErr != nil {
			errResponse := util.HandleErrors(serviceErr)
			c.JSON(errResponse.StatusCode, gin.H{"error": errResponse.Message, "success": errResponse.Success})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "deposit transaction successful", "success": true})
	}
}

func (s *TransactionController) Withdraw() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		account := dtos.WithdrawRequestDto{}

		if err := c.BindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
			return
		}

		_, serviceErr := s.transactionService.Withdraw(account)

		if serviceErr != nil {
			errResponse := util.HandleErrors(serviceErr)
			c.JSON(errResponse.StatusCode, gin.H{"error": errResponse.Message, "success": errResponse.Success})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "withdrawal transaction successful", "success": true})
	}
}

func (s *TransactionController) Transfer() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
