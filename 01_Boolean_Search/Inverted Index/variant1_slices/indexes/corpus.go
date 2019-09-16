package indexes

import (
	"fmt"
	"sort"
)

type Corpus struct {
	Indexes []InvertedIndex
}

func NewCorpus(docs []string) *Corpus {

	var indexes []InvertedIndex

	for i, data := range docs {
		dictionary := NewTermsFromSlice(SplitRaw(data))
		// remove duplicates
		//dictionary = dictionary.Filter(func (e string) bool {
		//	if !dictionary.ContainsManyTimes(e) {
		//		return true
		//	}
		//	return false
		//})
		for _, v := range dictionary {

			indexes = append(indexes, InvertedIndex{
				Term:  v,
				Documents: []string {fmt.Sprintf("Doc%d", i+1)},
			})
		}
	}
	c := &Corpus{Indexes: SortIndexes(indexes)}
	c.prepareCollection()
	return c
}

func (c *Corpus) prepareCollection()  {

	counter := 0

	// find the same terms in different documents and update posting lists
	for counter <= len(c.Indexes) {
		c.findRepetition(counter)
		counter++
	}

	c.sortPostingLists()

}

func (c *Corpus) findRepetition(counter int) {

	ok := true

	for ok {
		if counter+1 < len(c.Indexes) && c.Indexes[counter].Term == c.Indexes[counter+1].Term {
			c.Indexes[counter].Documents = append(c.Indexes[counter].Documents, c.Indexes[counter+1].Documents...)
			// remove next item
			c.Indexes = append(c.Indexes[:counter+1], c.Indexes[counter+2:]...)
		} else {
			ok = false
		}
	}

}

func (c *Corpus) sortPostingLists() {
	for _, v := range c.Indexes {
		sort.Strings(v.Documents)
	}
}