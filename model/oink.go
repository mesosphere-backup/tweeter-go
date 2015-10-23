package model

import (
	"time"
)

type Oink struct {
	ID           string
	Content      string
	CreationTime time.Time
	Handle       string
}
