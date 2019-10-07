package storage

import (
	"../spimi"
	"../corpus"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

const (
	outputFile = "blocks/index.dat"
	tempBlockSize = 5000
	termsInBlock = 4
)

func InitStorage(inputDir string) *corpus.BlockTree {

	var bt *corpus.BlockTree

	if !fileExists(outputFile) {
		bt = spimi.Spimi(inputDir, outputFile, tempBlockSize, termsInBlock)
	}

	bt = loadBTree(outputFile)

	return bt

}

func DeserializeBlock(path string) *corpus.SerializedCorpus{
	f , err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Println(err)
	}

	data := make([]byte, stat.Size())

	for {
		_, err = f.Read(data)
		if err != nil {
			if err == io.EOF {
				break // end of the file
			} else {
				fmt.Println("Error reading file", err);
				os.Exit(1)
			}
		}
	}

	return corpus.SerializedCorpusFromBlock(string(data))

}

func DeserializeTerm(term, path string) corpus.SerializedToken {
	return DeserializeBlock(path).Filter(func(token corpus.SerializedToken) bool{
		return token.Term == term
	}).Tokens[0]
}

func fileExists(path string) bool {
	// detect if file exists
	var _, err = os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return true

}

func loadBTree(path string) *corpus.BlockTree {

	f , err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Println(err)
	}

	data := make([]byte, stat.Size())

	for {
		_, err = f.Read(data)
		if err != nil {
			if err == io.EOF {
				break // end of the file
			} else {
				fmt.Println("Error reading file", err);
				os.Exit(1)
			}
		}
	}
	bt := corpus.BlockTreeFromGOB64(string(data))

	return bt

}

type TermRank struct {
	File string
	Score float32
}

//As a first step, we introduce the overlap score measure: the score of a document d is the
//sum, over all query terms, of the number of times each of the query terms
//occurs in d. We can refine this idea so that we add up not the number of
//occurrences of each query term t in d, but instead the tf-idf weight of each
//term in d.
// Intersect 2 terms and sort documents using their inverse document frequency
func Intersect(bt *corpus.BlockTree, term1, term2 string) []TermRank {

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