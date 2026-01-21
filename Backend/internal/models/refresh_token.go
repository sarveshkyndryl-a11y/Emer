package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Hash        string
	ExpiresAt  time.Time
	Revoked    bool
	ReplacedBy *uuid.UUID // 👈 THIS IS THE KEY
	CreatedAt  time.Time
}
