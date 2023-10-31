package session

import (
	"SimpleReader/server/models/user"
	"crypto/rsa"
	"time"
)

type Session struct {
	UID         string
	LastConnect time.Time
	IP          string
	Login       int             `json:"-"`
	User        *user.User      `json:"-"`
	Key         *rsa.PrivateKey `json:"-"`
}

func NewSession(uid, ip string) *Session {
	s := Session{}
	s.UID = uid
	s.LastConnect = time.Now().UTC()
	s.IP = ip
	s.Login = -5
	s.User = user.NewUser()
	return &s
}

func (s *Session) SetUser(u *user.User) {
	s.User = u
}

func (s *Session) IsLogin() bool {
	return s.Login == 1
}
