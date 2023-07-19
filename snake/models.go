package snake

import (
	"time"
)

type Change struct {
	ID        int64     `json:"id"`
	Revision  int       `json:"revision"`
	Message   string    `json:"message"`
	Filename  string    `json:"filename"`
	Code      []byte    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
