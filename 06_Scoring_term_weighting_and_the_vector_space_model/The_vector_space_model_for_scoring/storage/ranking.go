package storage

import (
	"math"
	"sort"
	"strings"
	"../corpus"
)

type TermRank struct {
	File string
	Score float32
}

//As a first step, we introduce the overlap score measure: the score of a document d is the
//sum, over all query terms, of the number of times each of the query terms
//occurs in d. We can refine this idea so that we add up not the number of
//occurrences of each query term t in d, but instead the tf-idf weight of each
//term in d.
// ITFScore 2 terms and sort documents using their inverse document frequency
func ITFScore(bt *corpus.BlockTree, term1, term2 string) []TermRank {

	res := make([]TermRank, 0)

	block1, ok1 := bt.Get(term1)
	block2, ok2 := bt.Get(term2)
	if !ok1 || !ok2 {
		return res
	}

	p1 := DeserializeTerm(term1, block1.(string))
	p2 := DeserializeTerm(term2, block2.(string))

	len1 := len(p1.Docs)
	len2 := len(p2.Docs)
	i, j := 0, 0

	for i != len1  && j != len2 {
		doc1 := p1.Docs[i]
		doc2 := p2.Docs[j]

		//   if docID(p1[i]) == docID(p2[j]):
		if doc1.DocID == doc2.DocID {
			res = append(res, TermRank{
				doc1.File,
				score(doc1, doc2, p1, p2),
			})
			i++
			j++
		} else if doc1.DocID < doc2.DocID {
			i++
		} else {
			j++
		}
	}

	return sortScores(res)

}

//Score(q, d) = ∑ tf-idf(t,d)
//             t∈q
func score(doc1, doc2 corpus.SerializedDoc, p1, p2 corpus.SerializedToken) float32 {

	sum := float32(0)

	sum += float32(p1.TotalFrequency)*doc1.InverseDocumentFrequency
	sum += float32(p2.TotalFrequency)*doc2.InverseDocumentFrequency

	return sum

}

// RankSorter sorts indexes by term score.
type RankSorter []TermRank

func (a RankSorter) Len() int           { return len(a) }
func (a RankSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RankSorter) Less(i, j int) bool { return a[i].Score > a[j].Score }

func sortScores(scores []TermRank) []TermRank {
	sort.Sort(RankSorter(scores))
	return scores
}

// Get top K results for givent query using cosine score and vector model
func CosineScore(bt *corpus.BlockTree, query string, top int) []TermRank {

	tokens := parseTokens(query)
	scores := make([]TermRank, 0)

	for _, t := range tokens {

		block, ok := bt.Get(t.Term)
		if !ok {
			continue
		}

		p := DeserializeTerm(t.Term, block.(string))
		for _, d := range p.Docs {

			if doc, ok := bt.Documents.Get(d.DocID); ok {
				document := doc.(*corpus.DocumentIndex)

				if frequency, ok := document.Get(t.Term); ok {
					ntf := frequency.(float32)
					doc := InputToken {
						Term:                        t.Term,
						NormalizedDocumentFrequency: ntf,
						InverseDocumentFrequency:    d.InverseDocumentFrequency,
						TFxIDF:                      ntf * d.InverseDocumentFrequency,
					}
					scores = append(scores, TermRank{
						File:  d.File,
						Score: CosineSimilarity(t, doc),
					})
				}

			}

		}

	}

	return getTopKResults(sortScores(scores), top)

}

//COSINESCORE(q)
//1 float Scores[N] = 0
//2 Initialize Length[N]
//3 for each query term t
//4 do calculate wt,q and fetch postings list for t
//5 	for each pair(d, tft,d) in postings list
//6 	do Scores[d] += wft,d × wt,q
//7 Read the array Length[d]
//8 for each d
//9 do Scores[d] = Scores[d]/Length[d]
//10 return Top K components of Scores[]

func getTopKResults(scores []TermRank, k int) []TermRank {

	if len(scores) <= k {
		return scores
	}

	return scores[:k]

}

func parseTokens(query string) []InputToken {

	//TODO: normalize query input
	input := strings.Split(query, " ")

	temp := map[string] int{}
	for _, term := range input {
		temp[term] += 1
	}

	tokens := make([]InputToken, 0)
	for key, value := range temp {
		idf := corpus.CountInverseDocumentFrequency(len(input), value)
		ndf := float32(value)/float32(len(input))
		tokens = append(tokens, InputToken {
			Term:                        key,
			NormalizedDocumentFrequency: ndf,
			InverseDocumentFrequency:    idf,
			TFxIDF:                      ndf*idf,
		})
	}

	return tokens

}

type InputToken struct {
	Term string
	NormalizedDocumentFrequency float32
	InverseDocumentFrequency float32
	TFxIDF float32
}



// DotProduct of 2 vectors
func DotProduct(doc1, doc2 InputToken) float32 {

	return doc1.TFxIDF * doc2.TFxIDF

}


// Return euclidean length for the given document
func EuclideanLength(doc InputToken) float32 {

	return float32(math.Sqrt(math.Pow(float64(doc.TFxIDF), 2)))

}


// Calculate cosine similarity for 2 documents
func CosineSimilarity(doc1, doc2 InputToken) float32 {

	return DotProduct(doc1, doc2) / (EuclideanLength(doc1) * EuclideanLength(doc2))

}


// DotProduct of 2 vectors
func DotProduct2(doc1, doc2 []InputToken) float32 {

	sum := float32(0)

	for i := 0; i < int(math.Min(float64(len(doc1)), float64(len(doc2)))); i++ {
		sum += doc1[i].TFxIDF * doc2[i].TFxIDF
	}

	return sum

}


// Return euclidean length for the given document
func EuclideanLength2(docs []InputToken) float32 {

	res := float64(0)

	for _, doc := range docs {
		res += math.Pow(float64(doc.TFxIDF), 2)
	}

	return float32(math.Sqrt(res))

}


// Calculate cosine similarity for 2 documents
func CosineSimilarity2(doc1, doc2 []InputToken) float32 {

	return DotProduct2(doc1, doc2) / (EuclideanLength2(doc1) * EuclideanLength2(doc2))

}
