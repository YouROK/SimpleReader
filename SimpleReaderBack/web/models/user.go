package models

import (
	"html/template"
)

type User struct {
	Login     string              `json:"login,omitempty"`
	PassHash  string              `json:"pass_hash"`
	Email     string              `json:"email"`
	Role      int                 `json:"role,omitempty"`
	ReadBooks map[string]BookInfo `json:"read_books,omitempty"`
	Style     *Style              `json:"style,omitempty"`
}

type Style struct {
	PageSlideEffect int          `json:"page_slide_effect" default:"0"`
	TextSize        float64      `json:"text_size" default:"1"`
	TextIndent      float64      `json:"text_indent" default:"1"`
	ParagraphIndent float64      `json:"paragraph_indent" default:"1.5"`
	TextBright      int          `json:"text_bright" default:"0"`
	DayTheme        bool         `json:"day_theme" default:"true"`
	FontName        template.CSS `json:"font_name"`
}
