package main

import (
	. "./indexes"
	"fmt"
)

var docs = []string {
	"new home sales top forecast home",
	"home sales rise in july",
	"increase in home sales in july",
	"forecast july new home sales rise",
}

func main() {
	fmt.Println(NewCorpus(docs))
}
