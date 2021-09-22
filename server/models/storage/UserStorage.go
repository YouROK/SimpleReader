package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"SimpleReader/server/models/users"
)

type UserStorage struct {
	path      string
	usersList map[string]*users.User
}

func NewUserStorage(path string) *UserStorage {
	u := UserStorage{}
	u.path = path
	u.usersList = make(map[string]*users.User)
	return &u
}

func (us *UserStorage) Exit() {
	if us != nil && us.usersList != nil {
		for _, u := range us.usersList {
			if u.Name != "" {
				us.SaveUser(u)
			}
		}
	}
}

func (us *UserStorage) GetUser(name string) (*users.User, error) {
	lowername := strings.ToLower(name)
	if u, ok := us.usersList[lowername]; ok {
		return u, nil
	}
	u, err := us.ReadUser(lowername)
	if err == nil {
		us.usersList[name] = u
	}
	return u, err
}

func (us *UserStorage) ReadUser(name string) (*users.User, error) {
	userPath := path.Join(us.path, name, "info")
	data, err := ioutil.ReadFile(userPath)
	if err == nil {
		user := users.NewUser()
		err = json.Unmarshal(data, &user)
		if err == nil {
			return user, nil
		}
	}
	log.Println("Error read user", userPath, err)
	return nil, err
}

func (us *UserStorage) SaveUser(u *users.User) error {
	if u.Name == "" || u == nil {
		return errors.New("User is wrong")
	}
	userPath := path.Join(us.path, strings.ToLower(u.Name), "info")
	data, err := json.Marshal(u)
	if err != nil || len(data) == 0 {
		log.Println("Error save user", userPath, err)
	} else {
		os.MkdirAll(path.Join(us.path, strings.ToLower(u.Name)), 0777)
		err = ioutil.WriteFile(userPath, data, 0777)
		if err == nil {
			us.usersList[strings.ToLower(u.Name)] = u
		}
	}
	return err
}
