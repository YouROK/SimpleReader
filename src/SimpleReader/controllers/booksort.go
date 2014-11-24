package controllers

import (
	"SimpleReader/models/fbreader"
	"SimpleReader/models/storage"
	"SimpleReader/models/users"
	"sort"
)

type SortBooksBySequence []*fbreader.XMLTitleInfo

func (a SortBooksBySequence) Len() int      { return len(a) }
func (a SortBooksBySequence) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortBooksBySequence) Less(i, j int) bool {
	ai := a[i].Sequence.Name + a[i].Sequence.Number
	aj := a[j].Sequence.Name + a[j].Sequence.Number
	return ai < aj
}

type bookInfoSlice []*booksinfo

type booksinfo struct {
	BookDesc *fbreader.XMLTitleInfo
	BookInfo users.BookInfo
}

func (a bookInfoSlice) Len() int { return len(a) }
func (a bookInfoSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a bookInfoSlice) Less(i, j int) bool {
	if (a[i].BookInfo.LastReadPage == 0 && a[j].BookInfo.LastReadPage == 0) || (a[i].BookInfo.LastRead.UTC().UnixNano() == a[j].BookInfo.LastRead.UTC().UnixNano()) {
		ai := a[i].BookDesc.Author.FirstName + a[i].BookDesc.Author.LastName + a[i].BookDesc.Author.MiddleName + a[i].BookDesc.BookTitle + a[i].BookDesc.Sequence.Number + a[i].BookDesc.Sequence.Name
		aj := a[j].BookDesc.Author.FirstName + a[j].BookDesc.Author.LastName + a[j].BookDesc.Author.MiddleName + a[j].BookDesc.BookTitle + a[j].BookDesc.Sequence.Number + a[j].BookDesc.Sequence.Name
		return ai < aj
	} else {
		ai := a[i].BookInfo.LastRead
		aj := a[j].BookInfo.LastRead
		return ai.After(aj)
	}
}

func getSortedBooksByDate(userBooks map[string]users.BookInfo) []*fbreader.XMLTitleInfo {
	bookslist := make(bookInfoSlice, 0, len(userBooks))

	for _, v := range userBooks {
		hash := v.BookHash
		desc, err := storage.GetStorage().GetBookStorage().GetBookDesc(hash)
		if err == nil {
			desc.Hash = hash
			bookslist = append(bookslist, &booksinfo{desc, v})
		}
	}

	sort.Sort(bookslist)

	ret := make([]*fbreader.XMLTitleInfo, 0, len(bookslist))
	for _, v := range bookslist {
		ret = append(ret, v.BookDesc)
	}
	return ret
}
