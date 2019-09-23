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

	//home, ok1 := c.Get("home")
	//sales, ok2 := c.Get("sales")
	//ok := ok1 && ok2
	//
	//if ok {
	//	fmt.Println(c.Intersect(home.(corpus.Index), sales.(corpus.Index)))
	//	fmt.Println(c.PositionalIntersect(home.(corpus.Index), sales.(corpus.Index), 2))
	//}

	//fo*st  $fo AND st$ -> forecast, forest
	fmt.Println(c.KGrammTermsIntersect("$fo", "st$"))
	//sa*es  $sa AND es$ -> sales, sanderes
	fmt.Println(c.KGrammTermsIntersect("$sa", "es$"))

	fmt.Println(c.KGrammTermsIntersect("$re", "red"))

}