package models

import (
	"time"
)

type BookInfo struct {
	BookHash        string
	LastReadSection int
	LastRead        time.Time
}
