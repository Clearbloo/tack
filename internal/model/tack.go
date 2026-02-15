package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type TackKind string

const (
	KindPin  TackKind = "pin"
	KindTodo TackKind = "todo"
	KindWarn TackKind = "warn"
)

type Tack struct {
	ID        string     `json:"id"`
	Kind      TackKind   `json:"kind"`
	Message   string     `json:"message"`
	Dir       string     `json:"-"`
	Done      bool       `json:"done,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	DoneAt    *time.Time `json:"done_at,omitempty"`
}

// NewID generates a short 4-character hex ID.
func NewID() string {
	b := make([]byte, 2)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Emoji returns the display emoji for this tack kind.
func (k TackKind) Emoji() string {
	switch k {
	case KindPin:
		return "📌"
	case KindTodo:
		return "☐"
	case KindWarn:
		return "⚠️"
	default:
		return "•"
	}
}

// DoneEmoji returns the emoji for a completed todo.
func DoneEmoji() string {
	return "☑"
}
