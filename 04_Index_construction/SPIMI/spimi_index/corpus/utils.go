package corpus

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"strings"
)

type SerializedToken struct {
	Term string
	Docs []SerializedDoc
}

type SerializedDoc struct {
	Positions []int
	DocID int
	File string
	Frequency int
}

type SerializedCorpus struct {
	Tokens []SerializedToken
}


// Go binary encoder
func (corpus *Corpus) ToGOB64() string {

	// convert to serialized corpus
	tokens := make([]SerializedToken, 0)
	corpus.Each(func(key, value interface{}) {
		term := key.(string)
		docs := value.(Index).Docs.Values()
		documents := make([]SerializedDoc, 0)
		for _, d := range docs {
			doc := d.(Doc)
			documents = append(documents, SerializedDoc{
				Positions: doc.Positions,
				DocID:     doc.ID,
				File:      doc.File,
				Frequency: doc.Frequency,
			})
		}
		tokens = append(tokens, SerializedToken{
			Term: term,
			Docs: documents,
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

	corpus := NewCorpus(2)
	corpus.BuildIndexFromSerializedTokens(sCorpus.Tokens)

	for _, token := range sCorpus.Tokens {
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
	}

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