package corpus

import (
	"fmt"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/sets/hashset"
)

type Index struct {
	Docs           Docs
	TotalFrequency int32
}

func (this Index) Contains(id int) bool {
	_, contains := this.Docs.Get(id)
	return contains
}

type KGrammIndex struct {
	*hashmap.Map
	k int
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

type SoundexTerms struct {
	*hashset.Set
}
