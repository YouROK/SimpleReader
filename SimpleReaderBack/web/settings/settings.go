package settings

import (
	"log"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

var (
	Path string
	db   *bolt.DB
)

func Open() {
	var err error
	db, err = bolt.Open(filepath.Join(Path, "sr.db"), 0666, nil)
	if err != nil {
		log.Panic(err)
	}
}

func Close() {
	if db != nil {
		db.Close()
		db = nil
	}
}

func name(name string) []byte {
	return []byte(name)
}
