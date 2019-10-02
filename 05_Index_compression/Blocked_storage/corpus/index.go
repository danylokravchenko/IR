package corpus

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/sets/hashset"
	"./automaton"
	"sync"
)

type Index struct {
	Docs           Docs
	TotalFrequency int
}

func (this Index) Contains(id int) bool {
	_, contains := this.Docs.Get(id)
	return contains
}


// Update document's frequency, position and append new document
func (index *Index) UpdateDocument(id int, positions []int) {

	document, _ := index.Docs.Get(id)
	doc := document.(Doc)
	doc.Frequency++
	doc.Positions = append(doc.Positions, positions...)

}

type KGrammIndex struct {
	*hashmap.Map
	k int
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}

func (kgramm *KGrammIndex) Print() {
	fmt.Println(kgramm.k, "GrammIndex")
	for _, v := range kgramm.Keys() {
		fmt.Printf("Key - %s, values - \n", v)
		terms_, _ := kgramm.Get(v)
		terms := terms_.(KGrammTerms)
		for _, t := range terms.Values() {
			fmt.Printf("%s, ", t)
		}
		fmt.Println()
	}
}

type KGrammTerms struct {
	*hashset.Set
}

type SoundexIndex struct {
	*hashmap.Map
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}

type SoundexTerms struct {
	*hashset.Set
}

type Token struct {
	Term string
	Position int
	DocID int
	File string
}

type Automaton struct {
	*automaton.Tree
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}

type BlockTree struct{
	*hashmap.Map
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
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(&SerializedBlockTree{blocks})
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

	bt := &BlockTree{hashmap.New()}

	for _, b := range sbt.Blocks {
		bt.Put(b.Term, b.Block)
	}

	return bt

}