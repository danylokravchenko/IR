package corpus

import (
	"fmt"
	"github.com/dotcypress/phonetics"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/hashset"
)

// Build inverted index from slice
func (corpus *Corpus) BuildIndexFromSlice(data []string) {

	corpus.wg.Add(len(data) * 4)

	for i, s := range data {
		words := splitRaw(s)
		go corpus.createIndex(words, i)
		go corpus.buildKGrammIndexFromTerms(words)
		go corpus.buildSoundexIndexFromTerms(words)
		go corpus.buildAutomatonIndexFromTerms(words)
	}

	corpus.wg.Wait()

}


// Build inverted index from parsed tokens
func (corpus *Corpus) BuildIndexFromTokens(tokens []Token) {

	corpus.wg.Add(len(tokens))

	for _, t := range tokens {
		go corpus.createIndexFromToken(t)
	}

	corpus.wg.Wait()

}

// Build inverted index from parsed tokens
func (corpus *Corpus) BuildIndexFromParsedTokens(tokens []Token) {

	corpus.wg.Add(len(tokens) * 4)

	for _, t := range tokens {
		go corpus.createIndexFromToken(t)
		go corpus.buildKGrammIndex(t.Term)
		go corpus.buildSoundexIndex(t.Term)
		go corpus.buildAutomatonIndex(t.Term)
	}

	corpus.wg.Wait()

}

// Build inverted index from parsed tokens
func (corpus *Corpus) BuildIndexFromSerializedTokens(tokens []SerializedToken) {

	corpus.wg.Add(len(tokens) * 4)

	for _, t := range tokens {
		go corpus.createIndexFromSerializedToken(t)
		go corpus.buildKGrammIndex(t.Term)
		go corpus.buildSoundexIndex(t.Term)
		go corpus.buildAutomatonIndex(t.Term)
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


// Create or update index for terms
func (corpus *Corpus) createIndexFromSerializedToken(token SerializedToken) {

	corpus.mutex.Lock()

	if index, ok := corpus.Get(token.Term); !ok {
		docs := treemap.NewWithIntComparator()
		totalFrequency := 0
		for _, d := range token.Docs {
			docs.Put(d.DocID, Doc{
				ID:        d.DocID,
				File:      d.File,
				Frequency: d.Frequency,
				Positions: d.Positions,
			})
			totalFrequency += d.Frequency
		}

		corpus.Put(token.Term, Index{Docs{docs}, totalFrequency})
	} else {
		documents := index.(Index)
		for _, d := range token.Docs {
			documents.TotalFrequency += d.Frequency
			if !documents.Contains(d.DocID) {
				documents.Docs.Put(d.DocID, Doc{
					ID:        d.DocID,
					File:      d.File,
					Frequency: d.Frequency,
					Positions: d.Positions,
				})
			} else {
				documents.UpdateDocument(d.DocID, d.Positions)
			}
		}
	}

	corpus.mutex.Unlock()

	corpus.wg.Done()

}


// Save kgramm keywords into map
func (corpus *Corpus) buildKGrammIndexFromTerms(terms []string) {

	for _, term := range terms {

		corpus.kGramm.mutex.Lock()

		gramms := splitKGramm(term, corpus.kGramm.k)

		for _, g := range gramms {

			if index, ok := corpus.kGramm.Get(g); !ok {
				corpus.kGramm.Put(g, KGrammTerms{hashset.New(term)})
			} else {
				terms := index.(KGrammTerms)
				terms.Add(term) //duplicates ignores
			}

		}

		corpus.kGramm.mutex.Unlock()
	}

	corpus.wg.Done()

}

// Save kgramm keywords into map
func (corpus *Corpus) buildKGrammIndex(term string) {

	corpus.kGramm.mutex.Lock()

	gramms := splitKGramm(term, corpus.kGramm.k)

	for _, g := range gramms {

		if index, ok := corpus.kGramm.Get(g); !ok {
			corpus.kGramm.Put(g, KGrammTerms{hashset.New(term)})
		} else {
			terms := index.(KGrammTerms)
			terms.Add(term) //duplicates ignores
		}

	}

	corpus.wg.Done()

	corpus.kGramm.mutex.Unlock()

}

// Save soundex value into map (English)
func (corpus *Corpus) buildSoundexIndexFromTerms(terms []string) {

	for _, term := range terms {

		corpus.soundex.mutex.Lock()

		val := phonetics.EncodeSoundex(term)

		if index, ok := corpus.soundex.Get(val); !ok {
			corpus.soundex.Put(val, SoundexTerms{hashset.New(term)})
		} else {
			terms := index.(SoundexTerms)
			terms.Add(term) //duplicates ignores
		}

		corpus.soundex.mutex.Unlock()

	}

	corpus.wg.Done()

}


// Save soundex value into map (English)
func (corpus *Corpus) buildSoundexIndex(term string) {

	corpus.soundex.mutex.Lock()

	val := phonetics.EncodeSoundex(term)

	if index, ok := corpus.soundex.Get(val); !ok {
		corpus.soundex.Put(val, SoundexTerms{hashset.New(term)})
	} else {
		terms := index.(SoundexTerms)
		terms.Add(term) //duplicates ignores
	}

	corpus.soundex.mutex.Unlock()


	corpus.wg.Done()

}

// Build Levenshtein Sparse automaton indexes
func (corpus *Corpus) buildAutomatonIndexFromTerms(terms []string) {

	for _, term := range terms {

		if term == "" {
			corpus.wg.Done()
			return
		}

		corpus.automaton.mutex.Lock()

		corpus.automaton.Insert(term)

		corpus.automaton.mutex.Unlock()

	}

	corpus.wg.Done()

}

// Build Levenshtein Sparse automaton indexes
func (corpus *Corpus) buildAutomatonIndex(term string) {

	if term == "" {
		corpus.wg.Done()
		return
	}

	corpus.automaton.mutex.Lock()

	corpus.automaton.Insert(term)

	corpus.automaton.mutex.Unlock()

	corpus.wg.Done()

}


