package services

import (
	"banking_application/api/database"
	"banking_application/api/domain/dtos"
	"banking_application/api/util"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type IUserService interface {
	CreateUser(newUser dtos.UserVM) (bool, error)
	GetUsers() ([]dtos.UserVM, error)
	GetUser(userID string) (dtos.UserVM, error)
}

type UserService struct {
	dbContext *database.AppDBContext
}

func NewUserService(dbContext *database.AppDBContext) *UserService {
	return &UserService{dbContext: dbContext}
}

var validate = validator.New()

func (s *UserService) CreateUser(newUser dtos.UserDto) (bool, error) {
	// Implement the logic to create a user
	validationErr := validate.Struct(newUser)
	if validationErr != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return false, fmt.Errorf("validation %v", validationErr)
	}

	// Hash the default password
	hashedPassword, err := util.HashPassword(newUser.Password)
	if err != nil {
		return false, fmt.Errorf("password hashing failed")
	}

	// Call the CreateUser method of the user database/repository
	_, err = s.dbContext.Insert(`INSERT INTO users (firstname, lastname, email, password, created_at) 
	VALUES ($1, $2, $3, $4, $5)`, newUser.Firstname, newUser.Lastname, newUser.Email, hashedPassword, time.Now())
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *UserService) GetUsers() ([]dtos.UserVM, error) {
	// Execute the SQL query to fetch users
	rows, err := s.dbContext.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to store users
	var users []dtos.UserVM

	// Iterate over the rows
	for rows.Next() {
		// Create a new UserVM instance to store the current row data
		var user dtos.UserVM

		// Scan the row into the UserVM struct
		if err := rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}

		// Append the UserVM to the slice
		users = append(users, user)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetUser(uuid uuid.UUID) (dtos.UserVM, error) {
	// Execute the SQL query to fetch users
	row := s.dbContext.QueryRow(`SELECT firstname, lastname, email FROM users where id = $1`, uuid)

	var user dtos.UserVM

	// Scan the row into the UserVM struct
	if err := row.Scan(&user.Firstname, &user.Lastname, &user.Email); err != nil {
		// Check if no rows were returned
		if errors.Is(err, sql.ErrNoRows) {
			return dtos.UserVM{}, fmt.Errorf("user not found")
		}
		return dtos.UserVM{}, err
	}

	return user, nil
}
