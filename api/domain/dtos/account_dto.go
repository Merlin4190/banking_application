package dtos

import (
	"github.com/google/uuid"
	"time"
)

type AccountDto struct {
	ID             string
	UserId         string
	AccountNumber  string
	AccountBalance float32
	CheckSum       *string
	IsActive       bool
	UpdatedAt      time.Time
}

type OpenAccountDto struct {
	UserId          uuid.UUID `json:"user_id" validate:"required,uuid"`
	DepositedAmount float32   `json:"amount" validate:"required"`
	IsActive        bool
}
