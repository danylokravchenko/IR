package corpus

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
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

func (sc *SerializedCorpus) ToGOB64() string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(sc)
	if err != nil { fmt.Println(`failed gob Encode`, err) }

	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func SerializedCorpusFromBlock(str string) *SerializedCorpus{
	sc := &SerializedCorpus{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil { fmt.Println(`failed base64 Decode`, err); }
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(sc)
	if err != nil { fmt.Println(`failed gob Decode`, err); }

	return sc
}

func (this *SerializedCorpus) Filter(filter func(token SerializedToken) bool) *SerializedCorpus {
	new_ := make([]SerializedToken, 0, len(this.Tokens))
	for _, v := range this.Tokens {
		if filter(v) {
			new_ = append(new_, v)
		}
	}
	return &SerializedCorpus{ Tokens: new_ }
}


type SerializedBlockTree struct {
	Blocks []SerializedBlock
}

type SerializedBlock struct {
	Term string
	Block string
}