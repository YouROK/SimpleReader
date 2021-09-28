package settings

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	bolt "go.etcd.io/bbolt"
)

func init() {
	go func() {
		ldir := filepath.Join(Path, "links")
		os.MkdirAll(ldir, 0777)
		for {
			checkLinks()
			files, err := ioutil.ReadDir(ldir)
			if err == nil {
				for _, ff := range files {
					if !ff.IsDir() {
						addLink(filepath.Join(ldir, ff.Name()))
					}
				}
			}
			time.Sleep(30 * time.Second)
		}
	}()
}

type RegLink struct {
	EMail   string
	Hash    string
	Expired time.Time
}

func addLink(path string) {
	buf, err := ioutil.ReadFile(path)
	if err == nil {
		rl := new(RegLink)
		rl.EMail = string(buf)
		rl.Hash = genHash()
		rl.Expired = time.Now().Add(time.Hour * 24 * 2)

		db.Update(func(tx *bolt.Tx) error {
			links, err := tx.CreateBucketIfNotExists(name("Links"))
			if err != nil {
				return err
			}
			buf, err := json.Marshal(rl)
			if err == nil {
				links.Put(name(rl.Hash), buf)
			}
			return err
		})
	}
}

func checkLinks() {
	db.Update(func(tx *bolt.Tx) error {
		links, err := tx.CreateBucketIfNotExists(name("Links"))
		if err != nil {
			return err
		}
		rems := make([][]byte, 0)
		links.ForEach(func(hash, buf []byte) error {
			var rl *RegLink
			if json.Unmarshal(buf, &rl) == nil {
				if time.Now().After(rl.Expired) {
					rems = append(rems, hash)
				}
			}
			return nil
		})
		for _, rem := range rems {
			links.Delete(rem)
		}
		return err
	})
}

func GetLink(hash string) *RegLink {
	var rl *RegLink
	db.View(func(tx *bolt.Tx) error {
		links, err := tx.CreateBucketIfNotExists(name("Links"))
		if err != nil {
			return err
		}

		buf := links.Get(name(hash))
		return json.Unmarshal(buf, &rl)
	})
	return rl
}

func RemLink(hash string) {
	db.Update(func(tx *bolt.Tx) error {
		links, err := tx.CreateBucketIfNotExists(name("Links"))
		if err != nil {
			return err
		}
		links.Delete(name(hash))
		return nil
	})
}

func genHash() string {
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	hash := sha512.Sum512([]byte(ts))
	return hex.EncodeToString(hash[:])
}
