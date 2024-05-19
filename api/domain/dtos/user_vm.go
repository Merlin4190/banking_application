package dtos

import (
	"time"

	"github.com/google/uuid"
)

type UserVM struct {
	ID        uuid.UUID
	Firstname string
	Lastname  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
