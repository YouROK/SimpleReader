package storage

import (
	"SimpleReader/models/global"
	"SimpleReader/models/sessions"
	"crypto/sha1"
	"encoding/hex"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"time"
)

type SessionStorage struct {
	path    string
	sesList map[string]*sessions.Session
}

func NewSessionStorage(path string) *SessionStorage {
	s := SessionStorage{}
	s.path = path
	s.sesList = make(map[string]*sessions.Session)
	return &s
}

func (ss *SessionStorage) Exit() {

}

func (ss *SessionStorage) NewSesion(req *http.Request) *sessions.Session {
	uid := generateUIDSession(req.UserAgent())
	s := sessions.NewSession(uid, req.RemoteAddr)
	ss.sesList[uid] = s
	return s
}

func (ss *SessionStorage) GetSession(req *http.Request) *sessions.Session {
	cookieSession, err := req.Cookie("session")
	if err == nil {
		uid := cookieSession.Value
		if s, ok := ss.sesList[uid]; ok {
			s.LastConnect = time.Now().UTC()
			s.IP = req.RemoteAddr
			return s
		}
	}
	return nil
}

func Sessions() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request) {
		ses := GetStorage().GetSessionStorage().GetSession(req)
		if ses == nil {
			ses = GetStorage().GetSessionStorage().NewSesion(req)
		}
		cookie := http.Cookie{}
		cookie.Name = "session"
		cookie.Value = ses.UID
		cookie.Expires = time.Now().Add(time.Hour * 24 * 2)
		cookie.HttpOnly = true
		cookie.Path = "/"
		req.Header.Set("Cookie", "session="+ses.UID)
		http.SetCookie(res, &cookie)
	}
	return nil
}

func generateUIDSession(salt string) string {
	hash := sha1.New()
	hash.Write([]byte(strconv.Itoa(int(time.Now().UnixNano())) + salt))
	s := hex.EncodeToString(hash.Sum(nil))
	return s
}

func cleaner() {
	for !global.Stoping {
		time.Sleep(5 * time.Minute)
		remList := make(map[string]*sessions.Session)
		for k, s := range GetStorage().GetSessionStorage().sesList {
			if s.LastConnect.UTC().Add(4 * 24 * time.Hour).After(time.Now().UTC()) {
				if s.Login == 1 {
					GetStorage().GetUserStorage().SaveUser(s.User)
				}
				remList[k] = s
			}
		}
		if len(remList) > 0 {
			for k, _ := range remList {
				delete(GetStorage().GetSessionStorage().sesList, k)
			}
		}
	}
}
