package settings

import (
	"encoding/json"
	"strings"

	bolt "go.etcd.io/bbolt"

	"SimpleReader/web/models"
)

func GetUser(login string) *models.User {
	var user *models.User
	db.View(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists(name("Users"))
		if err != nil {
			return err
		}
		buf := users.Get(name(strings.ToLower(login)))
		err = json.Unmarshal(buf, &user)
		if err != nil {
			return err
		}
		return nil
	})
	return user
}

func GetUsers() []*models.User {
	var usersList []*models.User
	db.View(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists(name("Users"))
		if err != nil {
			return err
		}
		users.ForEach(func(k, v []byte) error {
			var usr *models.User
			err := json.Unmarshal(v, &usr)
			if err == nil {
				usersList = append(usersList, usr)
			}
			return nil
		})
		return nil
	})
	return usersList
}

func AddUser(usr *models.User) {
	db.Update(func(tx *bolt.Tx) error {
		users, err := tx.CreateBucketIfNotExists(name("Users"))
		if err != nil {
			return err
		}
		buf, err := json.Marshal(usr)
		if err != nil {
			return err
		}
		users.Put(name(strings.ToLower(usr.Login)), buf)
		return nil
	})
}
