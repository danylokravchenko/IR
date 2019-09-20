package main

import (
	"fmt"
	//"github.com/aaaton/golem"
	//"github.com/aaaton/golem/dicts/en"
	"github.com/caneroj1/stemmer"
	"github.com/antoineaugusti/wordsegmentation"
	"github.com/antoineaugusti/wordsegmentation/corpus"
)

func main() {
	// the language packages are available under golem/dicts
	// "en" is for english
	//lemmatizer, err := golem.New(en.New())
	//if err != nil {
	//	panic(err)
	//}
	//word := lemmatizer.Lemma("Abducting")
	//if word != "abduct" {
	//	panic("The output is not what is expected!")
	//}

	str := "running"

	// stem a single word
	stem := stemmer.Stem(str)
	fmt.Println(stem)
	// stem = RUN

	strings := []string{
		"playing",
		"skies",
		"singed",
	}

	// stem a list of words
	stems := stemmer.StemMultiple(strings)
	fmt.Println(stems)
	// stems = [PLAI SKI SIN]

	// Grab the default English corpus that will be created thanks to TSV files
	englishCorpus := corpus.NewEnglishCorpus()
	fmt.Println(wordsegmentation.Segment(englishCorpus, "thisisatest"))

}
