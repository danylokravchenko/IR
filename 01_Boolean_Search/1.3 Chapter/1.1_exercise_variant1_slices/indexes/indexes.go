package indexes

import (
	"sort"
	"strings"
)

type Documents []string

type Terms []string

type InvertedIndex struct {
	Term string
	Documents Documents
}

func SplitRaw(s string) []string {
	return strings.Split(s, " ")
}

func NewTermsFromSlice(strings []string) Terms {
	var terms = Terms{}
	for _, s := range strings {
		terms = terms.Append(s)
	}
	return terms
}

func (this Terms) ContainsManyTimes(s string) bool {
	counter := 0
	for _, data := range this {
		if data == s {
			counter++
		}
	}
	if counter > 1 {
		return true
	}
	return false
}

func (this Terms) Contains(s string) bool {
	for _, data := range this {
		if data == s {
			return true
		}
	}
	return false
}

func (this Terms) Append(item string) Terms{
	if !this.Contains(item) {
		return append(this, item)
	}
	return this
}

func (this Terms) Filter(predicate func (s string) bool) Terms {
	var terms Terms
	for _, v := range this {
		if predicate(v) {
			terms = append(terms, v)
		}
	}
	return terms
}

// IndexSorter sorts indexes by term name.
type IndexSorter []InvertedIndex

func (a IndexSorter) Len() int           { return len(a) }
func (a IndexSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IndexSorter) Less(i, j int) bool { return a[i].Term < a[j].Term }

func SortIndexes(indexes []InvertedIndex) []InvertedIndex {
	sort.Sort(IndexSorter(indexes))
	return indexes
}

