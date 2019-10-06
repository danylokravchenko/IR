package corpus

import (
	"sort"
	"sync"
)

type Zones struct {
	title  *TitleIndex
	corpus *Corpus
	mutex *sync.Mutex
	wg *sync.WaitGroup
}

// zone weights
const (
	titleWeight = 0.4
	bodyWeight = 0.6
)


func NewZoneIndex() *Zones {
	return &Zones{
		title:  NewTitleIndex(),
		corpus: NewCorpus(3),
		mutex: &sync.Mutex{},
		wg: &sync.WaitGroup{},
	}
}

// Build inverted index from parsed tokens
func (zones *Zones) BuildZonesIndexFromTokens(tokens []Token, fileTokens []Token) {

	zones.wg.Add(2)

	go zones.buildCorpus(tokens)
	go zones.buildTitleIndex(fileTokens)

	zones.wg.Wait()

}

func (zones *Zones) buildCorpus(tokens []Token) {

	zones.corpus.BuildIndexFromTokens(tokens)
	zones.wg.Done()

}

type TermRank struct {
	File string
	Score float32
}

// RankSorter sorts indexes by term name.
type RankSorter []TermRank

func (a RankSorter) Len() int           { return len(a) }
func (a RankSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RankSorter) Less(i, j int) bool { return a[i].Score > a[j].Score }

func sortScores(scores []TermRank) []TermRank {
	sort.Sort(RankSorter(scores))
	return scores
}

func (zones *Zones) buildTitleIndex(tokens []Token) {

	zones.title.BuildIndexFromTokens(tokens)
	zones.wg.Done()

}


func (zones *Zones) ZoneScore(term1, term2 string) []TermRank {

	scores := make([]TermRank, 0)

	index1, ok1 := zones.corpus.Get(term1)
	index2, ok2 := zones.corpus.Get(term2)
	if !ok1 || !ok2 {
		return scores
	}

	p1 := index1.(Index)
	p2 := index2.(Index)

	len1 := p1.Docs.Size() + 1
	len2 := p2.Docs.Size() + 1
	i, j := 1, 1

	for i != len1  && j != len2 {
		var (
			doc1, doc2           Doc
			document1, document2 interface{}
			ok1, ok2             bool
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
			scores = append(scores, TermRank{
				doc1.File,
				zones.weightedZone(doc1.ID, term1, term2),
			})//scores[docID(p1)] ← WEIGHTEDZONE(p1, p2, g)
			i++
			j++
		} else if ok && doc1.ID < doc2.ID {
			i++
		} else {
			j++
		}
	}

	return sortScores(scores)

}

//ZONESCORE(q1, q2)  q1, q2 - terms
//1 	float scores[N] = [0]
//2 	constant g[ℓ]
//3 	p1 ← postings(q1)
//4 	p2 ← postings(q2)
//5 	// scores[] is an array with a score entry for each document, initialized to zero.
//6 	//p1 and p2 are initialized to point to the beginning of their respective postings.
//7 	//Assume g[] is initialized to the respective zone weights.
//8 	while p1 != NIL and p2 != NIL
//9 		do if docID(p1) = docID(p2)
//10 			then scores[docID(p1)] ← WEIGHTEDZONE(p1, p2, g)
//11 			p1 ← next(p1)
//12 			p2 ← next(p2)
//13 		else if docID(p1) < docID(p2)
//14 			then p1 ← next(p1)
//15 		else p2 ← next(p2)
//16 return scores

// Function WEIGHTEDZONE is assumed to compute the inner loop
// sum of zone ranks
// TODO: reorganize in WEIGHTEDZONE to search if term exists in body and in title in this exactly fileID
func (zone *Zones) weightedZone(docID int, term1, term2 string) float32 {

	var score float32

	idx1, ok1 := zone.title.Get(term1)
	idx2, ok2 := zone.title.Get(term2)

	if ok1 && ok2 {
		index1 := idx1.(Index)
		index2 := idx2.(Index)
		if index1.Contains(docID) && index2.Contains(docID) {
			score += titleWeight
		}
	}

	idx1, ok1 = zone.corpus.Get(term1)
	idx2, ok2 = zone.corpus.Get(term2)

	if ok1 && ok2 {
		index1 := idx1.(Index)
		index2 := idx2.(Index)
		if index1.Contains(docID) && index2.Contains(docID) {
			score += bodyWeight
		}
	}

	return score

}