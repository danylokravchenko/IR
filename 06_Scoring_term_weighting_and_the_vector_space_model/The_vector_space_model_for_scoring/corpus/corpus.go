package corpus

import (
	"./automaton"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/treemap"
	"sync"
)

type Corpus struct {
	*treemap.Map
	TermsNum  int
	DocsNum   int
	kGramm    *KGrammIndex
	soundex   *SoundexIndex
	automaton *Automaton
	Documents *DocumentTree
	mutex     *sync.Mutex
	wg        *sync.WaitGroup
}


// New instance of Corpus with initialized map, kGramm map, soundex map and syncs
func NewCorpus() *Corpus{
	return &Corpus{
		treemap.NewWithStringComparator(),
		0,
		0,
		&KGrammIndex{
			Map: hashmap.New(),
			k: 3,
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
		&DocumentTree{
			treemap.NewWithIntComparator(),
			&sync.Mutex{},
			&sync.WaitGroup{},
		},
		&sync.Mutex{},
		&sync.WaitGroup{},
	}
}