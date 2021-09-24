package settings

import (
	"strings"

	bolt "go.etcd.io/bbolt"
)

func GetBook(id string) []byte {
	var book []byte
	db.View(func(tx *bolt.Tx) error {
		books, err := tx.CreateBucketIfNotExists(name("Books"))
		if err != nil {
			return err
		}
		book = books.Get(name(strings.ToLower(id)))
		return nil
	})
	return book
}

func GetBooks() []string {
	var booksIds []string
	db.View(func(tx *bolt.Tx) error {
		books, err := tx.CreateBucketIfNotExists(name("Books"))
		if err != nil {
			return err
		}
		books.ForEach(func(k, v []byte) error {
			booksIds = append(booksIds, string(k))
			return nil
		})
		return nil
	})
	return booksIds
}

func AddBook(id string, book []byte) {
	db.Update(func(tx *bolt.Tx) error {
		books, err := tx.CreateBucketIfNotExists(name("Books"))
		if err != nil {
			return err
		}
		books.Put(name(strings.ToLower(id)), book)
		return nil
	})
}
