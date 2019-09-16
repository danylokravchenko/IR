package main

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"strings"
)

type Index struct {
	*treemap.Map
}

type Doc struct {
	file  string
}

type Docs struct {
	docs []Doc
	frequency int32
}

func (this Docs) Contains(s string) bool {
	for _, el := range this.docs {
		if el.file == s {
			return true
		}
	}
	return false
}

var docs = []string {
	"new home sales top forecast home",
	"home sales rise in july",
	"increase in home sales in july",
	"forecast july new home sales rise",
}

func (index *Index) BuildIndexFromSlice(data []string) {
	for i, s := range data {
		index.createIndex(s, i)
	}
}

func (index *Index) createIndex(s string, counter int) {
	words := splitRaw(s)
	counter++
	file := fmt.Sprintf("Doc%d", counter)
	for _, w := range words {
		if docs, ok := index.Get(w); !ok {
			index.Put(w, Docs{ []Doc{{file:file}}, 1})
		} else {
			if !docs.(Docs).Contains(file) {
				documents := docs.(Docs)
				documents.docs = append(documents.docs, Doc{file:file})
				documents.frequency++
				index.Put(w, documents)
			}
		}
	}
}

func splitRaw(s string) []string {
	return strings.Split(strings.Trim(s, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@"), " ")
}

func main() {
	// initialize map
	index := &Index{Map:treemap.NewWithStringComparator()}

	// build index
	index.BuildIndexFromSlice(docs)
	fmt.Println(index.Keys())
	fmt.Println(index.Values())

}
