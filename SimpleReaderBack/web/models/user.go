package models

import (
	"html/template"
)

type User struct {
	Login     string              `json:"login"`
	PassHash  string              `json:"pass_hash"`
	Email     string              `json:"email"`
	Role      int                 `json:"role"`
	ReadBooks map[string]BookInfo `json:"read_books"`
}

type Style struct {
	PageSlideEffect int
	TextSize        float64
	TextIndent      float64
	ParagraphIndent float64
	TextBright      int
	DayTheme        bool
	FontName        template.CSS
}
