package spimi

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
)

type SerializedToken struct {
	Term string
	Docs []SerializedDoc
}

type SerializedDoc struct {
	Positions []int
	DocID int
	File string
	Frequency int32
}

type SerializedCorpus struct {
	Tokens []SerializedToken
}


// TODO: try json encoding
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

	corpus := &Corpus {treemap.NewWithStringComparator()}

	for _, token := range sCorpus.Tokens {
		if index, ok := corpus.Get(token.Term); !ok {
			docs := treemap.NewWithIntComparator()
			totalFrequency := int32(0)
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
				if !index.(Index).Contains(d.DocID) {
					documents.Docs.Put(d.DocID, Doc{
						ID:        d.DocID,
						File:      d.File,
						Frequency: d.Frequency,
						Positions: d.Positions,
					})
				} else {
					documents.updateDocument(d.DocID, d.Positions)
				}
			}
		}
	}

	return corpus

}

type Corpus struct {
	*treemap.Map
}

type Doc struct {
	ID        int // int because map comparator is int type
	File      string
	Frequency int32
	Positions []int
}

type Docs struct {
	*treemap.Map
}

type Index struct {
	Docs           Docs //[]Doc
	TotalFrequency int32
}

func (this Index) Contains(id int) bool {
	_, contains := this.Docs.Get(id)
	return contains
}


func (index *Index) updateDocument(id int, positions []int) {
	document, _ := index.Docs.Get(id)
	doc := document.(Doc)
	doc.Frequency += int32(len(positions))
	doc.Positions = append(doc.Positions, positions...)
}
