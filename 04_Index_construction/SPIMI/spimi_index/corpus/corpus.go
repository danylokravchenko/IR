package corpus

import (
	"./automaton"
	"fmt"
	"github.com/dotcypress/phonetics"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/hashset"
	"sync"
)

type Corpus struct {
	*treemap.Map
	kGramm  *KGrammIndex
	soundex *hashmap.Map
	automaton *automaton.Tree
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
		},
		hashmap.New(),
		automaton.NewTree(),
		&sync.Mutex{},
		&sync.WaitGroup{},
	}
}


// Build inverted index from slice
func (corpus *Corpus) BuildIndexFromSlice(data []string) {

	corpus.wg.Add(len(data) * 4)

	for i, s := range data {
		words := splitRaw(s)
		go corpus.createIndex(words, i)
		go corpus.buildKGrammIndex(words)
		go corpus.buildSoundexIndex(words)
		go corpus.buildAutomatonIndex(words)
	}

	corpus.wg.Wait()

}


// Build inverted index from slice
func (corpus *Corpus) BuildIndexFromTokens(tokens []Token) {

	corpus.wg.Add(len(tokens) * 3)

	for _, t := range tokens {
		go corpus.createIndexFromToken(t)
		go corpus.buildKGrammIndexFromToken(t)
		go corpus.buildAutomatonIndexFromToken(t)
	}

	corpus.wg.Wait()

}


// Create or update index for terms
func (corpus *Corpus) createIndex(words []string, id int) {

	id++
	file := fmt.Sprintf("Doc%d", id)

	for position, w := range words {

		corpus.mutex.Lock()

		if index, ok := corpus.Get(w); !ok {
			docs := treemap.NewWithIntComparator()
			docs.Put(id, Doc{
				ID:        id,
				File:      file,
				Frequency: 1,
				Positions: []int{position + 1},
			})
			corpus.Put(w, Index{ Docs{docs}, 1})
		} else {
			documents := index.(Index)
			documents.TotalFrequency++

			if !documents.Contains(id) {
				documents.Docs.Put(id, Doc{
					ID:        id,
					File:      file,
					Frequency: 1,
					Positions: []int{position + 1},
				})
			} else {
				documents.UpdateDocument(id, []int{position + 1})
			}
		}

		corpus.mutex.Unlock()

	}

	corpus.wg.Done()

}


// Create or update index for terms
func (corpus *Corpus) createIndexFromToken(token Token) {

	corpus.mutex.Lock()

	if index, ok := corpus.Get(token.Term); !ok {
		docs := treemap.NewWithIntComparator()
		docs.Put(token.DocID, Doc{
			ID:        token.DocID,
			File:      token.File,
			Frequency: 1,
			Positions: []int{token.Position},
		})
		corpus.Put(token.Term, Index{ Docs{docs}, 1})
	} else {
		documents := index.(Index)
		documents.TotalFrequency++

		if !documents.Contains(token.DocID) {
			documents.Docs.Put(token.DocID, Doc{
				ID:        token.DocID,
				File:      token.File,
				Frequency: 1,
				Positions: []int{token.Position},
			})
		} else {
			documents.UpdateDocument(token.DocID, []int{token.Position})
		}
	}

	corpus.mutex.Unlock()

	corpus.wg.Done()

}


// Save kgramm keywords into map
func (corpus *Corpus) buildKGrammIndex(terms []string) {

	for _, term := range terms {

		corpus.mutex.Lock()

		gramms := splitKGramm(term, corpus.kGramm.k)

		for _, g := range gramms {

			if index, ok := corpus.kGramm.Get(g); !ok {
				corpus.kGramm.Put(g, KGrammTerms{hashset.New(term)})
			} else {
				terms := index.(KGrammTerms)
				terms.Add(term) //duplicates ignores
				// don't need next line (I hope :) )
				//corpus.kGramm.Put(g,terms)
			}

		}

		corpus.mutex.Unlock()
	}

	corpus.wg.Done()

}

// Save kgramm keywords into map
func (corpus *Corpus) buildKGrammIndexFromToken(token Token) {

	corpus.mutex.Lock()

	gramms := splitKGramm(token.Term, corpus.kGramm.k)

	for _, g := range gramms {

		if index, ok := corpus.kGramm.Get(g); !ok {
			corpus.kGramm.Put(g, KGrammTerms{hashset.New(token.Term)})
		} else {
			terms := index.(KGrammTerms)
			terms.Add(token.Term) //duplicates ignores
		}

	}

	corpus.wg.Done()

	corpus.mutex.Unlock()

}

// Save soundex value into map (English)
func (corpus *Corpus) buildSoundexIndex(terms []string) {

	for _, term := range terms {

		corpus.mutex.Lock()

		val := phonetics.EncodeSoundex(term)

		if index, ok := corpus.soundex.Get(val); !ok {
			corpus.soundex.Put(val, SoundexTerms{hashset.New(term)})
		} else {
			terms := index.(SoundexTerms)
			terms.Add(term) //duplicates ignores
		}

		corpus.mutex.Unlock()

	}

	corpus.wg.Done()

}


// Build Levenshtein Sparse automaton indexes
func (corpus *Corpus) buildAutomatonIndex(terms []string) {

	for _, term := range terms {

		corpus.mutex.Lock()

		corpus.automaton.Insert(term)

		corpus.mutex.Unlock()

	}

	corpus.wg.Done()

}

// Build Levenshtein Sparse automaton indexes
func (corpus *Corpus) buildAutomatonIndexFromToken(token Token) {

	if token.Term == "" {
		corpus.wg.Done()
		return
	}

	corpus.mutex.Lock()

	corpus.automaton.Insert(token.Term)

	corpus.mutex.Unlock()

	corpus.wg.Done()

}

