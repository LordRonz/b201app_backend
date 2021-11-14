package auth

import (
	"time"
)

// Token holds refresh jwt information.
type Token struct {
	ID        int       `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Token      string    `json:"-"`
	Expiry     time.Time `json:"-"`
	Mobile     bool      `sql:",notnull" json:"mobile"`
	Identifier string    `json:"identifier,omitempty"`
}
