package main

import (
	"./lemmatization"
	"fmt"
	"log"
)

func main() {
	lemmatazier, err := lemmatization.New("./dictionary/lemmatization-uk.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(lemmatazier.GetLema("аваріях"))
	fmt.Println(lemmatazier.Contains("аваріях"))

}
