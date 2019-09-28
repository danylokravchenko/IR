package corpus

import (
	"fmt"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/sets/hashset"
	"./automaton"
	"sync"
)

type Index struct {
	Docs           Docs
	TotalFrequency int
}

func (this Index) Contains(id int) bool {
	_, contains := this.Docs.Get(id)
	return contains
}


// Update document's frequency, position and append new document
func (index *Index) UpdateDocument(id int, positions []int) {

	document, _ := index.Docs.Get(id)
	doc := document.(Doc)
	doc.Frequency++
	doc.Positions = append(doc.Positions, positions...)

}

type KGrammIndex struct {
	*hashmap.Map
	k int
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}

func (kgramm *KGrammIndex) Print() {
	fmt.Println(kgramm.k, "GrammIndex")
	for _, v := range kgramm.Keys() {
		fmt.Printf("Key - %s, values - \n", v)
		terms_, _ := kgramm.Get(v)
		terms := terms_.(KGrammTerms)
		for _, t := range terms.Values() {
			fmt.Printf("%s, ", t)
		}
		fmt.Println()
	}
}

type KGrammTerms struct {
	*hashset.Set
}

type SoundexIndex struct {
	*hashmap.Map
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}

type SoundexTerms struct {
	*hashset.Set
}

type Token struct {
	Term string
	Position int
	DocID int
	File string
}

type Automaton struct {
	*automaton.Tree
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}
