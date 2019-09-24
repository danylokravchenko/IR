package main

import (
	"./corpus"
	"fmt"
)

var docs = []string {
	"new home sales top forecast home retired",
	"home sales rise in july june red redemption",
	"increase in home sales in july forest",
	"forecast july new home sales rise sanderes", //sanderes just to fit 'sa*es'
	"Miller Muller",
}

func main() {
	//fmt.Println(splitKGramm("castle",3)) //castle: $ca, ast, stl, tle, le$
	//fmt.Println(splitKGramm("dolphin",3)) //dolphin: $do, olp, lph, phi, hin, in$

	// initialize corpus
	c := corpus.New(3)

	// build index
	c.BuildIndexFromSlice(docs)

	// print corpus
	//c.Print()



	//fmt.Println(c.Intersect("home", "sales"))
	//fmt.Println(c.PositionalIntersect("home", "sales", 2))

	//fo*st  $fo AND st$ -> forecast, forest
	//fmt.Println(c.KGrammTermsIntersect("$fo", "st$"))
	//sa*es  $sa AND es$ -> sales, sanderes
	//fmt.Println(c.KGrammTermsIntersect("$sa", "es$"))

	//fmt.Println(c.KGrammTermsIntersect("$re", "red"))

	fmt.Println(c.GetSimilarlySoundWords("Miller")) // Miller, Muller

}