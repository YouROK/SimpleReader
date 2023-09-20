package user

import (
	"html/template"
	"time"
)

type BookInfo struct {
	BookHash     string
	LastReadPage int
	LastRead     time.Time
	SetPage      bool `json:"-"`
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

type User struct {
	Name         string
	Pass         string
	Email        string
	Avatar       string
	Role         int
	ReadingBooks map[string]BookInfo
	Style        *Style
}

func NewUser() *User {
	usr := &User{}
	usr.ReadingBooks = make(map[string]BookInfo)
	usr.Style = NewStyle()
	return usr
}

func (u *User) GetBookInfo(bookHash string) BookInfo {
	bookInfo, ok := u.ReadingBooks[bookHash]
	if !ok {
		bookInfo = BookInfo{BookHash: bookHash, LastRead: time.Now().UTC()}
	}
	return bookInfo
}

func (u *User) SetBookInfo(binfo BookInfo) {
	u.ReadingBooks[binfo.BookHash] = binfo
}

func (u *User) RemoveBookInfo(hash string) {
	delete(u.ReadingBooks, hash)
}

func NewStyle() *Style {
	s := &Style{}
	s.PageSlideEffect = 0
	s.TextSize = 1.0
	s.TextIndent = 1.0
	s.ParagraphIndent = 1.5
	s.TextBright = 0
	s.DayTheme = true
	return s
}

func (s *Style) GetTheme() string {
	if s.DayTheme {
		return "d"
	}
	return "n"
}

func (s *Style) GetColor() int {
	if s.DayTheme {
		return s.TextBright
	}
	return 255 - s.TextBright
}

func (s *Style) GetBackGround() string {
	if s.DayTheme {
		return "#ffffff"
	}
	return "#000000"
}
