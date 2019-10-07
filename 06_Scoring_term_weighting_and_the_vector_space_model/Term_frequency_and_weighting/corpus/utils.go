package corpus

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"math"
	"strings"
)


// Go binary encoder
func (corpus *Corpus) ToGOB64() string {

	// convert to serialized corpus
	tokens := make([]SerializedToken, 0)
	corpus.Each(func(key, value interface{}) {
		// TODO: save Frequency here too
		term := key.(string)
		index := value.(Index)
		docs := index.Docs.Values()
		documents := make([]SerializedDoc, 0)
		for _, d := range docs {
			doc := d.(Doc)
			documents = append(documents, SerializedDoc{
				Positions: doc.Positions,
				DocID:     doc.ID,
				File:      doc.File,
				Frequency: doc.Frequency,
				InverseDocumentFrequency: doc.InverseDocumentFrequency,
			})
		}
		tokens = append(tokens, SerializedToken{
			Term: term,
			Docs: documents,
			TotalFrequency: index.TotalFrequency,
		})
	})
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(&SerializedCorpus{Tokens:tokens})
	if err != nil { fmt.Println(`failed gob Encode`, err) }

	return base64.StdEncoding.EncodeToString(b.Bytes())

}

// Go binary decoder
func FromGOB64(str string) *Corpus {

	sCorpus := &SerializedCorpus{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil { fmt.Println(`failed base64 Decode`, err); }
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(sCorpus)
	if err != nil { fmt.Println(`failed gob Decode`, err); }

	corpus := NewCorpus()
	corpus.BuildIndexFromSerializedTokens(sCorpus.Tokens)

	return corpus

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


//How is the document frequency df of a term used to scale its weight?
//Denoting as usual the total number of documents in a collection by N, we define
//the inverse document frequency (idf) of a term t as follows:
//idf(t) = log (N/df(t))
// Count inverse documents frequency for terms in the corpus
func (corpus *Corpus) CountInverseDocumentsFrequency() {

	corpus.Each(func(key, value interface{}) {
		index := value.(Index)
		invDocFreq := float32(math.Log(float64(corpus.DocsNum)/float64(index.TotalFrequency)))
		if invDocFreq <= 0 {
			index.InverseDocumentFrequency = 0.001
		} else {
			index.InverseDocumentFrequency = invDocFreq
		}

		index.Docs.Each(func(key, value interface{}) {
			doc := value.(Doc)
			invDocFreq := float32(math.Log(float64(corpus.DocsNum)/float64(doc.Frequency)))
			if invDocFreq <= 0 {
				doc.InverseDocumentFrequency = 0.001
			} else {
				doc.InverseDocumentFrequency = invDocFreq
			}
			index.Docs.Put(key, doc)
		})

	})

}