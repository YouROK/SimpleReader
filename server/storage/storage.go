package storage

import (
	"SimpleReader/server/models/session"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

var (
	storage *Storage
)

type Storage struct {
	path        string
	bookStorage *BookStorage
	userStorage *UserStorage
	sesStorage  *SessionStorage
}

func NewStorage(path string) {
	storage = &Storage{}
	storage.path = path
}

func GetStorage() *Storage {
	return storage
}

func GetBookStorage() *BookStorage {
	return GetStorage().GetBookStorage()
}

func GetSessionStorage() *SessionStorage {
	return GetStorage().GetSessionStorage()
}

func GetUserStorage() *UserStorage {
	return GetStorage().GetUserStorage()
}

func GetSession(c *gin.Context) *session.Session {
	return GetStorage().GetSessionStorage().GetSession(c)
}

func Sessions() gin.HandlerFunc {
	return func(c *gin.Context) {
		ses := GetStorage().GetSessionStorage().GetSession(c)
		if ses == nil {
			ses = GetStorage().GetSessionStorage().NewSesion(c)
		}
		cookie := http.Cookie{}
		cookie.Name = "session"
		cookie.Value = ses.UID
		cookie.HttpOnly = true
		cookie.MaxAge = 48 * 3600
		cookie.Path = "/"
		c.Request.Header.Set("Cookie", "session="+ses.UID)
		http.SetCookie(c.Writer, &cookie)
		c.Next()
	}
}

func (s *Storage) GetBookStorage() *BookStorage {
	if s.bookStorage == nil {
		s.bookStorage = newBookStorage(path.Join(s.path, "books"))
	}
	return s.bookStorage
}

func (s *Storage) GetUserStorage() *UserStorage {
	if s.userStorage == nil {
		s.userStorage = newUserStorage(path.Join(s.path, "users"))
	}
	return s.userStorage
}

func (s *Storage) GetSessionStorage() *SessionStorage {
	if s.sesStorage == nil {
		s.sesStorage = newSessionStorage()
	}
	return s.sesStorage
}

func (s *Storage) GetPath() string {
	return s.path
}

func (s *Storage) Exit() {
	s.sesStorage.Exit()
	s.userStorage.Exit()
	s.bookStorage.Exit()
}
