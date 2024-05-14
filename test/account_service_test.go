package test

import (
	"banking_application/api/database"
	"banking_application/api/services"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"
)

// TestGenerateAccountNumber tests the GenerateAccountNumber method
func TestGenerateAccountNumber(t *testing.T) {
	accountNumber := "5213434566756"
	// Mock the DBContext
	mockDB := database.MockDBContext{}

	accountService := services.NewAccountService(&mockDB)

	// Test case: valid account number generation
	length := 10
	startingDigit := 5
	mockDB.On("QueryRow", mock.Anything, mock.AnythingOfType("string")).Return(new(sql.Row)).Once()
	mockDB.On("QueryRow", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()
	accountNumber, err := accountService.GenerateAccountNumber(length, startingDigit)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(accountNumber) != length {
		t.Errorf("expected account number length %d, got %d", length, len(accountNumber))
	}
	if accountNumber[0] != byte(startingDigit+'0') {
		t.Errorf("expected starting digit %d, got %c", startingDigit, accountNumber[0])
	}

	// Test case: account number already exists in the database
	mockDB.On("QueryRow", mock.Anything, mock.AnythingOfType("string")).Return(new(sql.Row)).Once()
	_, err = accountService.GenerateAccountNumber(length, startingDigit)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows error, got: %v", err)
	}

	// Test case: invalid length (less than 1)
	invalidLength := 0
	_, err = accountService.GenerateAccountNumber(invalidLength, startingDigit)
	if err == nil {
		t.Error("expected error for invalid length, got nil")
	}

	// Test case: invalid starting digit (less than 0 or greater than 9)
	invalidStartingDigit := -1
	_, err = accountService.GenerateAccountNumber(length, invalidStartingDigit)
	if err == nil {
		t.Error("expected error for invalid starting digit, got nil")
	}
	invalidStartingDigit = 10
	_, err = accountService.GenerateAccountNumber(length, invalidStartingDigit)
	if err == nil {
		t.Error("expected error for invalid starting digit, got nil")
	}

	// Assert that all expected methods were called
	mockDB.AssertExpectations(t)
}
