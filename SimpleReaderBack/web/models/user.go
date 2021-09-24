package models

type User struct {
	Login     string              `json:"login"`
	PassHash  string              `json:"pass_hash"`
	Email     string              `json:"email"`
	Role      int                 `json:"role"`
	ReadBooks map[string]BookInfo `json:"read_books"`
}
