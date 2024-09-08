package security

import (
	"time"
)

// Maker is an interface  for managing tokens
type Maker interface {
	CreateToken(userId int, duration time.Duration, access bool) (string, *Payload, error)
	VerifiyToken(token string) (*Payload, error)
}
