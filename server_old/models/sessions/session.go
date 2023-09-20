package sessions

import (
	"crypto/rsa"
	"time"

	"SimpleReader/server_old/models/users"
)

type Session struct {
	UID         string
	LastConnect time.Time
	IP          string
	Login       int             `json:"-"`
	User        *users.User     `json:"-"`
	Key         *rsa.PrivateKey `json:"-"`
}

func NewSession(uid, ip string) *Session {
	s := Session{}
	s.UID = uid
	s.LastConnect = time.Now().UTC()
	s.IP = ip
	s.Login = -5
	s.User = users.NewUser()
	return &s
}

func (s *Session) SetUser(u *users.User) {
	s.User = nil
	s.User = u
}
