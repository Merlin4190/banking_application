package controllers

import (
	"banking_application/api/domain/dtos"
	"banking_application/api/services"
	"banking_application/api/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AccountController struct {
	accountService services.Account
}

func NewAccountController(service services.Account) *AccountController {
	return &AccountController{
		accountService: service,
	}
}

func (s *AccountController) CreateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var account dtos.OpenAccountDto

		if err := c.BindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
			return
		}

		_, serviceErr := s.accountService.OpenNewAccount(account)

		if serviceErr != nil {
			errResponse := util.HandleErrors(serviceErr)
			c.JSON(errResponse.StatusCode, gin.H{"error": errResponse.Message, "success": errResponse.Success})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Account created successfully", "success": true})

	}
}

func (s *AccountController) GetAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := s.accountService.GetAccounts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching accounts"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func (s *AccountController) GetAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		accountNumber := c.Param("account_number")
		result, err := s.accountService.GetAccount(accountNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching user account details"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func (s *AccountController) DeactivateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		accountNumber := c.Param("account_number")
		result, err := s.accountService.DeactivateAccount(accountNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while deactivating user account details"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
