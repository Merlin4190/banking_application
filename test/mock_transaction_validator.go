package test

import (
	"banking_application/api/domain/dtos"
	"github.com/stretchr/testify/mock"
)

type MockValidator struct {
	mock.Mock
}

func (v *MockValidator) IsTransactionReferenceExist(transactionReference string) (bool, error) {
	args := v.Called(transactionReference)
	return args.Bool(0), args.Error(1)
}

func (v *MockValidator) ValidateChecksum(account dtos.AccountDto) (bool, error) {
	args := v.Called(account)
	return args.Bool(0), args.Error(1)
}

func (v *MockValidator) ComputeChecksum(account dtos.AccountDto) (string, error) {
	args := v.Called(account)
	return args.String(0), args.Error(1)
}

func (v *MockValidator) NotANegativeNumber(number float32) bool {
	args := v.Called(number)
	return args.Bool(0)
}
