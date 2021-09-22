package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"SimpleReader/server/models/fbreader"
	"SimpleReader/server/models/sessions"
	"SimpleReader/server/models/storage"
	"SimpleReader/server/utils"
)

func (l *Library) AuthorsPage(stor storage.Storages, req *http.Request, r render.Render) {
	bookshash := stor.GetBookStorage().GetBooks()
	alphmap := make(map[string]string)
	var alph []string

	for _, h := range bookshash {
		b, e := stor.GetBookStorage().GetBookDesc(h)
		if e == nil && b.Author.LastName != "" {
			a := string([]rune(b.Author.LastName)[0])
			alphmap[a] = a
		}
	}
	alph = append(alph, "*")
	for k, _ := range alphmap {
		alph = append(alph, k)
	}
	sort.Strings(alph)

	type AuthorPage struct {
		Session  *sessions.Session
		Alphabet []string
	}

	r.HTML(200, "library/authorsPage", AuthorPage{stor.GetSessionStorage().GetSession(req), alph})
}

type Author struct {
	FirstName  string
	MiddleName string
	LastName   string
	Hash       string
}
type Authors []Author

func (a Authors) Len() int      { return len(a) }
func (a Authors) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Authors) Less(i, j int) bool {
	ai := a[i].LastName + a[i].FirstName + a[i].MiddleName
	aj := a[j].LastName + a[j].FirstName + a[j].MiddleName
	return ai < aj
}
func (a *Author) GetHash() string {
	a.Hash = utils.Md5HashStr(a.LastName + a.FirstName + a.MiddleName)
	return a.Hash
}

func HasPrefixInCase(s, prefix string) bool {
	s = strings.ToLower(s)
	prefix = strings.ToLower(prefix)
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func (l *Library) AuthorsSearch(stor storage.Storages, req *http.Request, r render.Render) {

	type AuthorSearch struct {
		Sort string
	}

	authorSearch := AuthorSearch{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&authorSearch)
	if err != nil {
		r.JSON(400, nil)
		return
	}
	log.Println(authorSearch)
	bookshash := stor.GetBookStorage().GetBooks()
	authorsmap := make(map[Author]string)

	for _, h := range bookshash {
		b, e := stor.GetBookStorage().GetBookDesc(h)
		if e == nil && b.Author.LastName != "" {
			a := string([]rune(b.Author.LastName)[0])
			if a == authorSearch.Sort || authorSearch.Sort == "*" {
				author := Author{b.Author.FirstName, b.Author.MiddleName, b.Author.LastName, ""}
				author.GetHash()
				authorsmap[author] = ""
			}
		}
	}

	authors := []Author{}
	for a, _ := range authorsmap {
		authors = append(authors, a)
	}
	sort.Sort(Authors(authors))
	r.JSON(200, authors)
}

func (l *Library) AuthorPage(stor storage.Storages, req *http.Request, r render.Render, params martini.Params) {
	authorParam := params["author"]
	if authorParam == "" {
		r.Error(400)
		return
	}
	bookshash := stor.GetBookStorage().GetBooks()
	authorbooks := []*fbreader.XMLTitleInfo{}

	for _, h := range bookshash {
		b, e := stor.GetBookStorage().GetBookDesc(h)
		if e == nil {
			authorTmp := Author{b.Author.FirstName, b.Author.MiddleName, b.Author.LastName, ""}
			if authorParam == authorTmp.GetHash() {
				authorbooks = append(authorbooks, b)
			}
		}
	}

	if len(authorbooks) == 0 {
		r.Error(404)
		log.Println(authorParam)
		return
	}

	type AuthorPage struct {
		Author  Author
		Books   []*fbreader.XMLTitleInfo
		Session *sessions.Session
	}
	author := Author{authorbooks[0].Author.FirstName, authorbooks[0].Author.MiddleName, authorbooks[0].Author.LastName, ""}
	sort.Sort(SortBooksBySequence(authorbooks))
	r.HTML(200, "library/authorPage", AuthorPage{author, authorbooks, stor.GetSessionStorage().GetSession(req)})
}
