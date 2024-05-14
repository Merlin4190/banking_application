package dtos

import (
	"github.com/google/uuid"
	"time"
)

type AccountDto struct {
	ID             uuid.UUID
	UserId         uuid.UUID
	AccountNumber  string
	AccountBalance float32
	CheckSum       string
	IsActive       bool
	UpdatedAt      time.Time
}

type NewAccountDto struct {
	UserId          uuid.UUID `json:"userId" validate:"required,uuid"`
	DepositedAmount float32   `json:"balance" validate:"required"`
	IsActive        bool
}
