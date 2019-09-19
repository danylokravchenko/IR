package corpus

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"math"
	"strings"
)

type Corpus struct {
	*treemap.Map
}


func (corpus *Corpus) BuildIndexFromSlice(data []string) {
	for i, s := range data {
		corpus.createIndex(s, i)
	}
}

func (corpus *Corpus) createIndex(line string, id int) {
	words := splitRaw(line)
	id++
	file := fmt.Sprintf("Doc%d", id)
	for position, w := range words {
		if index, ok := corpus.Get(w); !ok {
			docs := treemap.NewWithIntComparator()
			docs.Put(id, Doc{
				ID:        id,
				File:      file,
				Frequency: 1,
				Positions: []int32{int32(position) + 1},
			})
			corpus.Put(w, Index{ Docs{docs}, 1})
		} else {
			documents := index.(Index)
			documents.TotalFrequency++
			if !index.(Index).Contains(id) {
				documents.Docs.Put(id, Doc{
					ID:        id,
					File:      file,
					Frequency: 1,
					Positions: []int32{int32(position) + 1},
				})
			} else {
				documents.updateDocument(id, position + 1)
			}
			corpus.Put(w, documents)
		}
	}
}

func (index *Index) updateDocument(id, position int) {
	document, _ := index.Docs.Get(id)
	doc := document.(Doc)
	doc.Frequency++
	doc.Positions = append(doc.Positions, int32(position))
	index.Docs.Put(id, doc)
}

func splitRaw(s string) []string {
	return strings.Split(strings.Trim(s, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@"), " ")
}

func (corpus *Corpus) Print() {
	corpus.Each(func(key interface{}, value interface{}) {
		index := value.(Index)
		fmt.Printf("term: %s, total Frequency: %d, posting list: \n",key.(string), index.TotalFrequency)
		index.Docs.Each(func(key interface{}, value interface{}) {
			fmt.Println(value.(Doc))
		})
	})
}

func (corpus *Corpus) PositionalIntersect(p1, p2 Index,  k int) Docs {
	var answer = Docs{treemap.NewWithIntComparator()}
	len1 := p1.Docs.Size() + 1
	len2 := p2.Docs.Size() + 1
	i, j := 1, 1
	for i != len1  && j != len2 {
		var(
			doc1, doc2 Doc
			document1, document2 interface{}
			ok1, ok2 bool
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
			var l []int32  // l <- ()
			pp1 := doc1.Positions
			pp2 := doc2.Positions

			plen1 := len(pp1)
			plen2 := len(pp2)
			ii, jj := 0, 0
			for ii != plen1 {
				for jj != plen2 {
					if math.Abs(float64(pp1[ii] - pp2[jj])) <= float64(k) {
						l = append(l, pp2[jj])
					} else if pp2[jj] > pp1[ii] {
						break
					}
					jj++
				}
				for len(l) > 0 && math.Abs(float64(l[0] - pp1[ii])) > float64(k){
					l = append(l[:0], l[1:]...)  // delete(l[0])
				}
				for _, ps := range l {
					answer.Put(doc1.ID, Doc {  // add answer(docID(p1), pos(pp1), ps)
						ID:        doc1.ID,
						File:      doc1.File,
						Frequency: 1,
						Positions: []int32{pp1[ii], ps},
					} )
				}
				ii++
			}
			i++
			j++
		} else if ok && doc1.ID < doc2.ID {
			i++
		} else {
			j++
		}
	}
	return answer
}
