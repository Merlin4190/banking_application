package dtos

import (
	"github.com/google/uuid"
	"time"
)

type AccountVM struct {
	ID             uuid.UUID
	UserId         uuid.UUID
	AccountNumber  string
	AccountBalance float32
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
