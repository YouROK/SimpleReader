package pages

import (
	"SimpleReader/server/api/utils"
	"SimpleReader/server/models/fb2"
	"SimpleReader/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

func AuthorPage(c *gin.Context) {
	queryHash := c.Param("hash")
	if queryHash == "" {
		c.Status(http.StatusNotFound)
		return
	}

	booksHash := storage.GetBookStorage().GetBooks()
	doneParse := make(chan *fb2.XMLTitleInfo)
	var infos []*fb2.XMLTitleInfo

	go func() {
		utils.ParallelFor(0, len(booksHash), 5, func(i int) {
			hash := booksHash[i]
			desc, err := storage.GetBookStorage().GetBookDesc(hash)
			if err == nil {
				authorTmp := &Author{desc.Author.FirstName, desc.Author.MiddleName, desc.Author.LastName, "", 0}
				if queryHash == authorTmp.GetHash() {
					doneParse <- desc
				}
			}
		})
		close(doneParse)
	}()

	for a := range doneParse {
		infos = append(infos, a)
	}

	if len(infos) == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	type AuthorPage struct {
		Author Author
		Books  []*fb2.XMLTitleInfo
	}

	sort.Slice(infos, func(i, j int) bool {
		ai := infos[i].Sequence.Name + infos[i].Sequence.Number
		aj := infos[j].Sequence.Name + infos[j].Sequence.Number
		return ai < aj
	})

	author := Author{infos[0].Author.FirstName, infos[0].Author.MiddleName, infos[0].Author.LastName, "", 0}
	c.HTML(http.StatusOK, "author.gohtml", AuthorPage{author, infos})
}
