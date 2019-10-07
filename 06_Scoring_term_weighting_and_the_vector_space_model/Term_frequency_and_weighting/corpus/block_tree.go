package corpus

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/emirpasic/gods/maps/hashmap"
)

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
