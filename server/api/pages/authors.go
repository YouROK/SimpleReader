package pages

import (
	"SimpleReader/server/api/utils"
	"SimpleReader/server/storage"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sort"
	"strings"
)

func SearchAuthorPage(c *gin.Context) {
	booksHash := storage.GetBookStorage().GetBooks()
	doneParse := make(chan string)
	alphmap := make(map[string]struct{})

	go func() {
		utils.ParallelFor(0, len(booksHash), 5, func(i int) {
			hash := booksHash[i]
			desc, err := storage.GetBookStorage().GetBookDesc(hash)
			if err == nil {
				aName := strings.TrimSpace(desc.Author.LastName)
				if aName != "" {
					a := string([]rune(aName)[0])
					doneParse <- strings.ToUpper(a)
				}
			}
		})
		close(doneParse)
	}()

	for a := range doneParse {
		alphmap[a] = struct{}{}
	}

	alph := append([]string{}, "*")
	for k, _ := range alphmap {
		alph = append(alph, k)
	}
	sort.Strings(alph)

	c.HTML(http.StatusOK, "authors.gohtml", alph)
}

type Author struct {
	FirstName  string
	MiddleName string
	LastName   string
	Hash       string
	Books      int
}
type Authors []Author

func (a *Author) GetHash() string {
	a.Hash = utils.Md5HashStr(a.LastName + a.FirstName + a.MiddleName)
	return a.Hash
}

func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}

func SearchAuthors(c *gin.Context) {
	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	query := strings.ToLower(string(buf))

	booksHash := storage.GetBookStorage().GetBooks()
	authorsmap := make(map[Author]int)
	doneParse := make(chan Author)

	go func() {
		utils.ParallelFor(0, len(booksHash), 5, func(i int) {
			hash := booksHash[i]
			desc, err := storage.GetBookStorage().GetBookDesc(hash)
			if err == nil && desc.Author.LastName != "" {
				if hasPrefix(desc.Author.LastName, query) || query == "*" {
					author := Author{desc.Author.FirstName, desc.Author.MiddleName, desc.Author.LastName, "", 1}
					author.GetHash()
					doneParse <- author
				}
			}
		})
		close(doneParse)
	}()

	for a := range doneParse {
		if _, ok := authorsmap[a]; ok {
			authorsmap[a]++
		} else {
			authorsmap[a] = 1
		}
	}

	var authors []Author
	for a, b := range authorsmap {
		a.Books = b
		authors = append(authors, a)
	}
	sort.Slice(authors, func(i, j int) bool {
		ai := authors[i].LastName + authors[i].FirstName + authors[i].MiddleName
		aj := authors[j].LastName + authors[j].FirstName + authors[j].MiddleName
		return ai < aj
	})
	c.JSON(200, authors)
}
