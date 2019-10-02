package corpus

import (
	"./automaton"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/treemap"
	"sync"
)

type Corpus struct {
	*treemap.Map
	kGramm  *KGrammIndex
	soundex *SoundexIndex
	automaton *Automaton
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}


// New instance of Corpus with initialized map, kGramm map, soundex map and syncs
func NewCorpus(kgrammSize int) *Corpus{
	return &Corpus{
		treemap.NewWithStringComparator(),
		&KGrammIndex{
			Map: hashmap.New(),
			k:kgrammSize,
			mutex: &sync.Mutex{},
			wg: &sync.WaitGroup{},
		},
		&SoundexIndex{
			hashmap.New(),
			&sync.Mutex{},
			&sync.WaitGroup{},
		},
		&Automaton{
			automaton.NewTree(),
			&sync.Mutex{},
			&sync.WaitGroup{},
		},
		&sync.Mutex{},
		&sync.WaitGroup{},
	}
}