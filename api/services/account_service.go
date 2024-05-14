package services

import (
	"banking_application/api/database"
	"banking_application/api/domain/dtos"
	"banking_application/api/util"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
)

type IAccountService interface {
	OpenNewAccount(account dtos.NewAccountDto) (bool, error)
	GetAccounts() ([]dtos.AccountVM, error)
	GetAccount(accountNumber string) (dtos.AccountVM, error)
	DeactivateAccount(accountNumber string) (bool, error)
}

type AccountService struct {
	dbContext database.DBContext
}

func NewAccountService(dbContext database.DBContext) *AccountService {
	return &AccountService{dbContext: dbContext}
}

func (s *AccountService) OpenNewAccount(request dtos.NewAccountDto) (bool, error) {
	validationErr := validate.Struct(request)

	if validationErr != nil {
		return false, fmt.Errorf("validation %v", validationErr)
	}

	//Generate New Account Number For User
	accountNumber, err := s.GenerateAccountNumber(10, 0)
	if err != nil {
		return false, err
	}

	_, insertErr := s.dbContext.Insert(`INSERT INTO accounts (user_id, account_number, account_balance, is_active) 
	VALUES ($1, $2, $3, $4)`, request.UserId, accountNumber, request.DepositedAmount, request.IsActive)

	if insertErr != nil {
		return false, err
	}
	return true, nil
}

func (s *AccountService) GetAccount(accountNumber string) (dtos.AccountVM, error) {
	row := s.dbContext.QueryRow(`SELECT account_number, account_balance, user_id FROM accounts where account_number = $1`, accountNumber)

	var accountDetails dtos.AccountVM

	// Scan the row into the UserVM struct
	if err := row.Scan(&accountDetails.AccountNumber, &accountDetails.AccountBalance, &accountDetails.UserId); err != nil {
		// Check if no rows were returned
		if errors.Is(err, sql.ErrNoRows) {
			return dtos.AccountVM{}, fmt.Errorf("account not found for accountnumber %s", accountNumber)
		}
		return dtos.AccountVM{}, err
	}
	return accountDetails, nil
}

func (s *AccountService) GetAccounts() ([]dtos.AccountVM, error) {
	rows, err := s.dbContext.Query(`SELECT * FROM accounts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to store accounts
	var accounts []dtos.AccountVM

	// Iterate over the rows
	for rows.Next() {
		// Create a new UserVM instance to store the current row data
		var accountDetails dtos.AccountVM

		// Scan the row into the AccountVM struct
		if err := rows.Scan(&accountDetails.ID, &accountDetails.AccountNumber, &accountDetails.AccountBalance, &accountDetails.IsActive,
			&accountDetails.CreatedAt, &accountDetails.UpdatedAt); err != nil {
			return nil, err
		}

		// Append the AccountVM to the slice
		accounts = append(accounts, accountDetails)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *AccountService) DeactivateAccount(accountNumber string) (bool, error) {
	var id string
	if err := s.dbContext.QueryRow(`SELECT id FROM accounts WHERE account_number = $1`, accountNumber).Scan(&id); err != nil {
		// Check if no rows were returned
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("account not found for accountnumber %s", accountNumber)
		}
		return false, err
	}

	//Update Account
	_, err := s.dbContext.Update(`UPDATE accounts SET is_active = false WHERE id = $1`, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *AccountService) GenerateAccountNumber(length int, startingDigit int) (string, error) {
	// Calculate min and max values for the random number
	min := util.PowerOf10(length - 1)
	max := util.PowerOf10(length) - 1

	// Generate a random number within the specified range
	randomNumber := rand.Intn(max-min+1) + min

	// Construct the account number by concatenating the starting digit and the random number
	accountNumber := fmt.Sprintf("%d%0*d", startingDigit, length-1, randomNumber)

	// Check if the generated account number already exists in the database
	var id string
	err := s.dbContext.QueryRow(`SELECT id FROM accounts WHERE account_number = $1`, accountNumber).Scan(&id)
	if err != nil {
		// Account number not found, return the generated account number
		if errors.Is(err, sql.ErrNoRows) {
			return accountNumber, nil
		}
		// Other error occurred
		return "", err
	}

	// Account number already exists, regenerate a new account number
	return s.GenerateAccountNumber(length, startingDigit)
}
