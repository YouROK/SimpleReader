package storage

import (
	"path"
	"sync"
)

type Storage struct {
	Storages
	path        string
	bookStorage *BookStorage
	userStorage *UserStorage
	sesStorage  *SessionStorage
	exitSync    sync.Mutex
}

func (s *Storage) SetPath(path string) {
	s.path = path
}

func (s *Storage) GetBookStorage() *BookStorage {
	s.exitSync.Lock()
	s.exitSync.Unlock()
	if s.bookStorage == nil {
		s.bookStorage = NewBookStorage(path.Join(s.path, "books"))
	}
	return s.bookStorage
}

func (s *Storage) GetUserStorage() *UserStorage {
	s.exitSync.Lock()
	s.exitSync.Unlock()
	if s.userStorage == nil {
		s.userStorage = NewUserStorage(path.Join(s.path, "users"))
	}
	return s.userStorage
}

func (s *Storage) GetSessionStorage() *SessionStorage {
	s.exitSync.Lock()
	s.exitSync.Unlock()
	if s.sesStorage == nil {
		s.sesStorage = NewSessionStorage(path.Join(s.path, "sessions"))
	}
	return s.sesStorage
}

var _init_ctx sync.Once
var _instance *Storage

func GetStorage() *Storage {
	_init_ctx.Do(func() {
		_instance = &Storage{}
	})
	return _instance
}

func (s *Storage) Exit() {
	s.exitSync.Lock()
	s.sesStorage.Exit()
	s.userStorage.Exit()
	s.bookStorage.Exit()
}

type Storages interface {
	GetBookStorage() *BookStorage
	GetUserStorage() *UserStorage
	GetSessionStorage() *SessionStorage
}
