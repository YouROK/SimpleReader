package storage

import (
	"SimpleReader/server/models/global"
	"SimpleReader/server/models/session"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type SessionStorage struct {
	sesList map[string]*session.Session
}

func newSessionStorage() *SessionStorage {
	s := SessionStorage{}
	s.sesList = make(map[string]*session.Session)
	return &s
}

func (ss *SessionStorage) Exit() {

}

func (ss *SessionStorage) NewSesion(c *gin.Context) *session.Session {
	uid := generateUIDSession(c.Request.UserAgent())
	s := session.NewSession(uid, c.RemoteIP())
	ss.sesList[uid] = s
	return s
}

func (ss *SessionStorage) GetSession(c *gin.Context) *session.Session {
	cookieSession, err := c.Request.Cookie("session")
	if err == nil {
		uid := cookieSession.Value
		if s, ok := ss.sesList[uid]; ok {
			s.LastConnect = time.Now().UTC()
			s.IP = c.RemoteIP()
			return s
		}
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
	for !global.Stoped {
		time.Sleep(5 * time.Minute)
		remList := make(map[string]*session.Session)
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
