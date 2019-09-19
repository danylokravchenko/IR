package main

import (
	"./corpus"
	"github.com/emirpasic/gods/maps/treemap"
)



var docs = []string {
	"new home sales top forecast home",
	"home sales rise in july",
	"increase in home sales in july",
	"forecast july new home sales rise",
}


func main() {
	// initialize corpus
	corpus := corpus.Corpus{treemap.NewWithStringComparator()}

	// build index
	corpus.BuildIndexFromSlice(docs)

	// print result
	corpus.Print()
}
