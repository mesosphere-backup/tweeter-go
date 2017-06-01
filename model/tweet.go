package model

import (
	"time"
)

type Tweet struct {
	ID string
	Content string
	CreationTime time.Time
	Handle string
}
