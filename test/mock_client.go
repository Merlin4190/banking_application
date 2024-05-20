package test

import (
	"banking_application/api/http/models"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (c *MockClient) PostPayment(url string, request models.PaymentRequest) (*models.PaymentResponse, error) {
	args := c.Called(url, request)
	response, _ := args.Get(0).(*models.PaymentResponse)
	return response, args.Error(1)
}

func (c *MockClient) GetPayment(url string) (*models.PaymentResponse, error) {
	args := c.Called(url)
	response, _ := args.Get(0).(*models.PaymentResponse)
	return response, args.Error(1)
}
