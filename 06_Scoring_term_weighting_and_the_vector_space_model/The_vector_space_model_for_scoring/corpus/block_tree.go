package corpus

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/treemap"
	"sync"
)

type BlockTree struct{
	*hashmap.Map
	Documents *DocumentTree
}

func (bt *BlockTree) ToGOB64() string {
	// convert to serialized block tree
	blocks := make([]SerializedBlock, 0)
	for _, key := range bt.Keys() {
		block, _ := bt.Get(key)
		blocks = append(blocks, SerializedBlock{
			Term:  key.(string),
			Block: block.(string),
		})
	}
	docs := make([]SerializedBlockDoc, 0)
	bt.Documents.Each(func(key, value interface{}) {
		doc := value.(DocumentIndex)
		terms := make([]SerializeBlockTerm, 0)
		doc.Each(func(key, value interface{}) {
			terms = append(terms, SerializeBlockTerm{
				Term:                key.(string),
				NormalizedFrequency: value.(float32),
			})
		})
		docs = append(docs, SerializedBlockDoc{
			DocID: key.(int),
			Terms: terms,
		})
	})
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(&SerializedBlockTree{blocks, docs})
	if err != nil { fmt.Println(`failed gob Encode`, err) }

	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func BlockTreeFromGOB64(str string) *BlockTree {
	sbt := &SerializedBlockTree{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil { fmt.Println(`failed base64 Decode`, err); }
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(sbt)
	if err != nil { fmt.Println(`failed gob Decode`, err); }

	bt := &BlockTree{
		hashmap.New(),
		&DocumentTree{
			treemap.NewWithIntComparator(),
			&sync.Mutex{},
			&sync.WaitGroup{},
		},
	}

	for _, b := range sbt.Blocks {
		bt.Put(b.Term, b.Block)
	}

	for _, d := range sbt.Documents {
		docs := &DocumentIndex{treemap.NewWithStringComparator()}
		for _, t := range d.Terms {
			docs.Put(t.Term, t.NormalizedFrequency)
		}
		bt.Documents.Put(d.DocID, docs)
	}

	return bt

}
