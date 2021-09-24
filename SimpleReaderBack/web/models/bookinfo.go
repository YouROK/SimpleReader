package models

import (
	"time"
)

type BookInfo struct {
	BookHash     string
	LastReadPage int
	LastRead     time.Time
}
