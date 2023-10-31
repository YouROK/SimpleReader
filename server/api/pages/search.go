package pages

import (
	"SimpleReader/server/api/utils"
	"SimpleReader/server/models/fb2"
	"SimpleReader/server/storage"
	"github.com/agnivade/levenshtein"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sort"
	"strings"
)

func SearchPage(c *gin.Context) {
	c.HTML(http.StatusOK, "search.gohtml", nil)
}

func SearchBooks(c *gin.Context) {
	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	query := strings.ToLower(string(buf))
	if len(query) < 3 {
		c.String(400, "Request is small")
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
				searchStr := strings.ToLower(desc.Author.FirstName + desc.Author.LastName + desc.Author.MiddleName + desc.BookTitle + desc.Sequence.Name + desc.Keywords)
				for _, v := range utils.GetGenres(desc.Genre) {
					searchStr += strings.ToLower(v)
				}
				if strings.Contains(searchStr, query) {
					doneParse <- desc
				}
			}
		})
		close(doneParse)
	}()

	for desc := range doneParse {
		infos = append(infos, desc)
	}

	sort.Slice(infos, func(i, j int) bool {
		lev1 := levenshtein.ComputeDistance(query, strings.ToLower(infos[i].BookTitle))
		lev2 := levenshtein.ComputeDistance(query, strings.ToLower(infos[j].BookTitle))
		return lev1 < lev2
	})

	html := ""

	for _, desc := range infos {
		html += `
<ons-card onclick="location.href='/book/` + desc.Hash + `'">
    <img style="float:left; margin-right: 10px; border-radius: 5%;" src="/cover/` + desc.Hash + `" width="86" height="130" />
    <div class="content">
        <div class="title"><b>` + desc.BookTitle + `</b><br/>
            ` + desc.Sequence.Name + ` ` + desc.Sequence.Number + `</div>
        <p>` + desc.Author.LastName + ` ` + desc.Author.MiddleName + ` ` + desc.Author.FirstName + `</p>
        <p><small><i>` + desc.GetGenresLine() + `</i></small></p>
        <div style="clear:both;"></div>
    </div>
</ons-card>`
	}

	c.Header("Content-Type", "text/plain; charset=UTF-8")
	c.String(200, html)
}
