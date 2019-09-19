package main

import (
	"./corpus"
	"fmt"
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
	c := corpus.Corpus{treemap.NewWithStringComparator()}

	// build index
	c.BuildIndexFromSlice(docs)

	// print corpus
	//c.Print()

	home, _ := c.Get("home")
	sales, _ := c.Get("sales")

	fmt.Println(c.PositionalIntersect(home.(corpus.Index), sales.(corpus.Index), 2))

}
