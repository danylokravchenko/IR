package corpus

import (
	"fmt"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/dotcypress/phonetics"
	"math"
	"strings"
	"sync"
)

type Corpus struct {
	*treemap.Map
	kGramm  *KGrammIndex
	soundex *hashmap.Map
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}


// New instance of Corpus with initialized map, kGramm map, soundex map and syncs
func New(kgrammSize int) *Corpus{
	return &Corpus{
		treemap.NewWithStringComparator(),
		&KGrammIndex{
			Map: hashmap.New(),
			k:kgrammSize,
		},
		hashmap.New(),
		&sync.Mutex{},
		&sync.WaitGroup{},
	}
}


// Build inverted index from slice
func (corpus *Corpus) BuildIndexFromSlice(data []string) {

	corpus.wg.Add(len(data))

	for i, s := range data {
		go corpus.createIndex(s, i)
	}

	corpus.wg.Wait()

}


// Create or update index for terms
func (corpus *Corpus) createIndex(line string, id int) {

	corpus.mutex.Lock()

	words := splitRaw(line)
	id++
	file := fmt.Sprintf("Doc%d", id)

	for position, w := range words {

		corpus.buildKGrammIndex(w)

		corpus.buildSoundexIndex(w)

		if index, ok := corpus.Get(w); !ok {
			docs := treemap.NewWithIntComparator()
			docs.Put(id, Doc{
				ID:        id,
				File:      file,
				Frequency: 1,
				Positions: []int32{int32(position) + 1},
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
					Positions: []int32{int32(position) + 1},
				})
			} else {
				documents.updateDocument(id, position + 1)
			}
			// don't need next line (I hope :) )
			//corpus.Put(w, documents)
		}

	}

	corpus.wg.Done()
	// defer is 40 ns/op
	corpus.mutex.Unlock()

}


// Save kgramm keywords into map
func (corpus *Corpus) buildKGrammIndex(term string) {

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

}


// Save soundex value into map (English)
func (corpus *Corpus) buildSoundexIndex(term string) {

	val := phonetics.EncodeSoundex(term)

	if index, ok := corpus.soundex.Get(val); !ok {
		corpus.soundex.Put(val, SoundexTerms{hashset.New(term)})
	} else {
		terms := index.(SoundexTerms)
		terms.Add(term) //duplicates ignores
	}

}


// Update document's frequency, position and append new document
func (index *Index) updateDocument(id, position int) {

	document, _ := index.Docs.Get(id)
	doc := document.(Doc)
	doc.Frequency++
	doc.Positions = append(doc.Positions, int32(position))
	// don't need next line (I hope :) )
	//index.Docs.Put(id, doc)

}

// Split line by 'space' and trim it
func splitRaw(s string) []string {
	return strings.Split(strings.Trim(s, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@"), " ")
}


// Build all available gramm for the given term
func splitKGramm(s string, k int) []string {

	var res []string
	l := len(s)

	if l < k {
		res = append(res, "$"+s)
		return res
	}

	if l == k {
		res = append(res, "$"+s[:k-1])
		res = append(res, s)
		res = append(res, s[1:] + "$")
		return res
	}

	res = append(res, "$"+s[:k-1])
	for i := 1; i <= l; i++ {
		if i+k+1 == l+k-1 {
			res = append(res, s[i:l] + "$")
			break
		}
		res = append(res, s[i:k+i])
	}

	return res

}


// Helper for printing all important information
func (corpus *Corpus) Print() {

	corpus.kGramm.Print()

	fmt.Println(corpus.soundex)

	corpus.Each(func(key interface{}, value interface{}) {
		index := value.(Index)
		fmt.Printf("term: %s, total Frequency: %d, posting list: \n",key.(string), index.TotalFrequency)
		index.Docs.Each(func(key interface{}, value interface{}) {
			fmt.Println(value.(Doc))
		})
	})

}


// Intersect kGramm Indexes for the given wildcard
func (corpus *Corpus) KGrammTermsIntersect(s1, s2 string) []string {

	var values1 []string
	var values2 []string
	var terms KGrammTerms

	if v1, ok1 := corpus.kGramm.Get(s1); ok1 {
		for _, v := range v1.(KGrammTerms).Values() {
			values1 = append(values1, v.(string))
		}
		terms = v1.(KGrammTerms)
	}

	if v2, ok2 := corpus.kGramm.Get(s2); ok2 {
		for _, v := range v2.(KGrammTerms).Values() {
			values2 = append(values2, v.(string))
		}
	}
	if s1 == "" {
		return values2
	}
	if s2 == "" {
		return values1
	}

	var res []string
	len1 := len(values1)
	len2 := len(values2)
	i, j := 0,0

	for i != len1 && j != len2 {
		if terms.Contains(values1[i]) && terms.Contains(values1[i]) {
			res = append(res, values1[i])
		}
		j++
		i++
	}

	return postFilter(res, s1, s2)

}


// Get terms that have the same soundex code
func (corpus *Corpus) GetSimilarlySoundWords(term string) []string {

	res := make([]string, 0)

	if terms, ok := corpus.soundex.Get(phonetics.EncodeSoundex(term)); ok {
		for _, term := range terms.(SoundexTerms).Values() {
			res = append(res, term.(string))
		}
	}

	return res

}

// Filter results to prevent terms with incorrect wildcards:
// red*  $re AND red -> retired !!! but it does not start with 'red'
func postFilter(terms []string, wildcard1, wildcard2 string) []string {

	var res []string

	for _, term := range terms {

		s1 := strings.Replace(wildcard1, "$", "", -1)
		s2 := strings.Replace(wildcard2, "$", "", -1)

		if strings.HasPrefix(wildcard1, "$") {
			if strings.HasPrefix(term, s1) && strings.Contains(term, s2) {
				res = append(res, term)
				continue
			}
		}

		if strings.HasSuffix(wildcard1, "$") {
			if strings.HasSuffix(term, s1) && strings.Contains(term, s2) {
				res = append(res, term)
				continue
			}
		}

		if strings.HasPrefix(wildcard2, "$") {
			if strings.HasPrefix(term, s2) && strings.Contains(term, s1) {
				res = append(res, term)
				continue
			}
		}

		if strings.HasSuffix(wildcard2, "$") {
			if strings.HasSuffix(term, s2) && strings.Contains(term, s1) {
				res = append(res, term)
				continue
			}
		}

		if strings.Contains(term, s1) && strings.Contains(term, s2) {
			res = append(res, term)
		}
	}

	return res

}


// Intersect Indexes by closest terms by their positions
func (corpus *Corpus) PositionalIntersect(term1, term2 string,  k int) Docs {

	index1, ok1 := corpus.Get(term1)
	index2, ok2 := corpus.Get(term2)
	if !ok1 || !ok2 {
		return Docs{}
	}

	p1 := index1.(Index)
	p2 := index2.(Index)

	var answer = Docs{treemap.NewWithIntComparator()}
	len1 := p1.Docs.Size() + 1
	len2 := p2.Docs.Size() + 1
	i, j := 1, 1

	for i != len1  && j != len2 {
		var(
			doc1, doc2 Doc
			document1, document2 interface{}
			ok1, ok2 bool
		)

		if document1, ok1 = p1.Docs.Get(i); ok1 { //check for nil
			doc1 = document1.(Doc)
		}

		if document2, ok2 = p2.Docs.Get(j); ok2 { //check for nil
			doc2 = document2.(Doc)
		}

		ok := ok1 && ok2
		//   if docID(p1[i]) == docID(p2[j]):
		if ok && doc1.ID == doc2.ID {
			var l []int32  // l <- ()
			pp1 := doc1.Positions
			pp2 := doc2.Positions

			plen1 := len(pp1)
			plen2 := len(pp2)
			ii, jj := 0, 0

			for ii != plen1 {
				for jj != plen2 {
					if math.Abs(float64(pp1[ii] - pp2[jj])) <= float64(k) {
						l = append(l, pp2[jj])
					} else if pp2[jj] > pp1[ii] {
						break
					}
					jj++
				}
				for len(l) > 0 && math.Abs(float64(l[0] - pp1[ii])) > float64(k){
					l = append(l[:0], l[1:]...)  // delete(l[0])
				}
				for _, ps := range l {
					answer.Put(doc1.ID, Doc {  // add answer(docID(p1), pos(pp1), ps)
						ID:        doc1.ID,
						File:      doc1.File,
						Frequency: 1,
						Positions: []int32{pp1[ii], ps},
					} )
				}
				ii++
			}
			i++
			j++
		} else if ok && doc1.ID < doc2.ID {
			i++
		} else {
			j++
		}
	}

	return answer

}


// Intersect 2 Indexes
func (corpus *Corpus) Intersect(term1, term2 string) Docs {

	index1, ok1 := corpus.Get(term1)
	index2, ok2 := corpus.Get(term2)
	if !ok1 || !ok2 {
		return Docs{}
	}

	p1 := index1.(Index)
	p2 := index2.(Index)

	var answer = Docs{treemap.NewWithIntComparator()}
	len1 := p1.Docs.Size() + 1
	len2 := p2.Docs.Size() + 1
	i, j := 1, 1

	for i != len1  && j != len2 {
		var(
			doc1, doc2 Doc
			document1, document2 interface{}
			ok1, ok2 bool
		)

		if document1, ok1 = p1.Docs.Get(i); ok1 { //check for nil
			doc1 = document1.(Doc)
		}

		if document2, ok2 = p2.Docs.Get(j); ok2 { //check for nil
			doc2 = document2.(Doc)
		}

		ok := ok1 && ok2
		//   if docID(p1[i]) == docID(p2[j]):
		if ok && doc1.ID == doc2.ID {
			answer.Put(doc1.ID, doc1)
			i++
			j++
		} else if ok && doc1.ID < doc2.ID {
			i++
		} else {
			j++
		}
	}
	return answer
}

//INTERSECT(p1, p2)
//1 answer ← ()
//2 while p1 != NIL and p2 != NIL
//3 do if docID(p1) = docID(p2)
//4 then ADD(answer, docID(p1))
//5 p1 ← next(p1)
//6 p2 ← next(p2)
//7 else if docID(p1) < docID(p2)
//8 then p1 ← next(p1)
//9 else p2 ← next(p2)
//10 return answer


//INTERSECTWITHSKIPS(p1, p2)
//1 answer ← ()
//2 while p1 != NIL and p2 != NIL
//3 do if docID(p1) = docID(p2)
//4 then ADD(answer, docID(p1))
//5 p1 ← next(p1)
//6 p2 ← next(p2)
//7 else if docID(p1) < docID(p2)
//8 then if hasSkip(p1) and (docID(skip(p1)) ≤ docID(p2))
//9 then while hasSkip(p1) and (docID(skip(p1)) ≤ docID(p2))
//10 do p1 ← skip(p1)
//11 else p1 ← next(p1)
//12 else if hasSkip(p2) and (docID(skip(p2)) ≤ docID(p1))
//13 then while hasSkip(p2) and (docID(skip(p2)) ≤ docID(p1))
//14 do p2 ← skip(p2)
//15 else p2 ← next(p2)
//16 return answer

//POSITIONALINTERSECT(p1, p2, k)
//1 answer ← ()
//2 while p1 != NIL and p2 != NIL
//3 do if docID(p1) = docID(p2)
//4 then l ← h i
//5 pp1 ← positions(p1)
//6 pp2 ← positions(p2)
//7 while pp1 6= NIL
//8 do while pp2 6= NIL
//9 do if |pos(pp1) − pos(pp2)| ≤ k
//10 then ADD(l, pos(pp2))
//11 else if pos(pp2) > pos(pp1)
//12 then break
//13 pp2 ← next(pp2)
//14 while l != () and |l[0] − pos(pp1)| > k
//15 do DELETE(l[0])
//16 for each ps ∈ l
//17 do ADD(answer,hdocID(p1), pos(pp1), psi)
//18 pp1 ← next(pp1)
//19 p1 ← next(p1)
//20 p2 ← next(p2)
//21 else if docID(p1) < docID(p2)
//22 then p1 ← next(p1)
//23 else p2 ← next(p2)
//24 return answer