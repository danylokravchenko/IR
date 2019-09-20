package main

import (
	"./corpus"
)

var docs = []string {
	"new home sales top forecast home",
	"home sales rise in july june",
	"increase in home sales in july",
	"forecast july new home sales rise",
}

func main() {
	//fmt.Println(splitKGramm("castle",3)) //castle: $ca, ast, stl, tle, le$
	//fmt.Println(splitKGramm("dolphin",3)) //dolphin: $do, olp, lph, phi, hin, in$

	// initialize corpus
	c := corpus.New(3)

	// build index
	c.BuildIndexFromSlice(docs)

	// print corpus
	c.Print()

	//home, _ := c.Get("home")
	//sales, _ := c.Get("sales")
	//
	//fmt.Println(c.Intersect(home.(corpus.Index), sales.(corpus.Index)))
	//
	//fmt.Println(c.PositionalIntersect(home.(corpus.Index), sales.(corpus.Index), 2))
}