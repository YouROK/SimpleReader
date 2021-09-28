package storage

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"

	"SimpleReader/web/models"
	"SimpleReader/web/settings"
)

func GetUser(mail string) *models.User {
	usrpath := path.Join(settings.Path, "users", strings.ToLower(mail), "info.json")
	buf, err := ioutil.ReadFile(usrpath)
	if err == nil {
		var usr *models.User
		json.Unmarshal(buf, usr)
		return usr
	}
	return nil
}

func GetUsers() []*models.User {
	usrpath := path.Join(settings.Path, "users")
	files, err := ioutil.ReadDir(usrpath)
	if err != nil {
		return nil
	}
	var users []*models.User
	for _, ff := range files {
		usr := GetUser(ff.Name())
		if usr != nil {
			users = append(users, usr)
		}
	}
	return users
}

func SetUser(usr *models.User) {
	usrpath := path.Join(settings.Path, "users", strings.ToLower(usr.Email), "info.json")
	buf, err := json.Marshal(usr)
	if err == nil {
		ioutil.WriteFile(usrpath, buf, 0664)
	}
}
