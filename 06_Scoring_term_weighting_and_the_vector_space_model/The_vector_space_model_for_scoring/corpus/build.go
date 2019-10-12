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

	corpus.wg.Add(len(tokens) * 2)

	for _, t := range tokens {
		go corpus.createIndexFromToken(t)
		go corpus.createDocumentIndexFromToken(t)
	}

	corpus.wg.Wait()

}

// Build inverted index from parsed tokens
func (corpus *Corpus) BuildIndexFromParsedTokens(tokens []Token) {

	corpus.wg.Add(len(tokens) * 5)

	for _, t := range tokens {
		go corpus.createIndexFromToken(t)
		go corpus.buildKGrammIndex(t.Term)
		go corpus.buildSoundexIndex(t.Term)
		go corpus.buildAutomatonIndex(t.Term)
		go corpus.createDocumentIndexFromToken(t)
	}

	corpus.wg.Wait()

}

// Build inverted index from parsed tokens
func (corpus *Corpus) BuildIndexFromSerializedTokens(tokens []SerializedToken) {

	corpus.wg.Add(len(tokens) * 5)

	for _, t := range tokens {
		go corpus.createIndexFromSerializedToken(t)
		go corpus.buildKGrammIndex(t.Term)
		go corpus.buildSoundexIndex(t.Term)
		go corpus.buildAutomatonIndex(t.Term)
		go corpus.createDocumentIndexFromSerializedToken(t)
	}

	corpus.wg.Wait()

}


// Create or update index for terms
func (corpus *Corpus) createIndex(words []string, id int) {

	id++
	file := fmt.Sprintf("Doc%d", id)
	corpus.DocsNum++

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
			corpus.Put(w, Index{ Docs{docs}, 1, 1, 0})
			corpus.TermsNum++
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
				documents.DocsNum++
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
		corpus.Put(token.Term, Index{ Docs{docs}, 1, 1, 0})
		corpus.TermsNum++
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
			documents.DocsNum++
		} else {
			documents.UpdateDocument(token.DocID, []int{token.Position})
		}

		corpus.Put(token.Term, documents)

	}

	corpus.mutex.Unlock()

	corpus.wg.Done()

}


// Create or update index for terms
func (corpus *Corpus) createIndexFromSerializedToken(token SerializedToken) {

	corpus.mutex.Lock()

	if index, ok := corpus.Get(token.Term); !ok {
		docs := treemap.NewWithIntComparator()
		for _, d := range token.Docs {
			docs.Put(d.DocID, Doc{
				ID:        d.DocID,
				File:      d.File,
				Frequency: d.Frequency,
				Positions: d.Positions,
			})
		}
		corpus.Put(token.Term, Index{Docs{docs}, token.TotalFrequency, 0, 0})
		corpus.TermsNum++
	} else {
		documents := index.(Index)
		documents.TotalFrequency += token.TotalFrequency
		for _, d := range token.Docs {
			if !documents.Contains(d.DocID) {
				documents.Docs.Put(d.DocID, Doc{
					ID:        d.DocID,
					File:      d.File,
					Frequency: d.Frequency,
					Positions: d.Positions,
				})
				documents.DocsNum++
			} else {
				documents.UpdateDocument(d.DocID, d.Positions)
			}
		}
		corpus.Put(token.Term, documents)
	}

	corpus.mutex.Unlock()

	corpus.wg.Done()

}


// Create or update document index for Documents
func (corpus *Corpus) createDocumentIndexFromToken(token Token) {

	corpus.Documents.mutex.Lock()

	if index, ok := corpus.Documents.Get(token.DocID); !ok {
		docs := DocumentIndex{
			Map:   treemap.NewWithStringComparator(),
		}
		docs.Put(token.Term, 1)
		corpus.Documents.Put(token.DocID, docs)
	} else {
		docs := index.(DocumentIndex)
		if d, ok := docs.Get(token.Term); !ok {
			docs.Put(token.Term, 1)
		} else {
			docs.Put(token.Term, d.(int)+1)
		}
	}

	corpus.Documents.mutex.Unlock()

	corpus.wg.Done()

}


// Create or update document index for Documents
func (corpus *Corpus) createDocumentIndexFromSerializedToken(token SerializedToken) {

	corpus.Documents.mutex.Lock()

	for _, doc := range token.Docs {

		if index, ok := corpus.Documents.Get(doc.DocID); !ok {
			docs := DocumentIndex{
				Map:   treemap.NewWithStringComparator(),
			}
			docs.Put(token.Term, token.TotalFrequency)
			corpus.Documents.Put(doc.DocID, docs)
		} else {
			docs := index.(DocumentIndex)
			if d, ok := docs.Get(token.Term); !ok {
				docs.Put(token.Term, token.TotalFrequency)
			} else {
				docs.Put(token.Term, d.(int)+token.TotalFrequency)
			}
		}

	}

	corpus.Documents.mutex.Unlock()

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


