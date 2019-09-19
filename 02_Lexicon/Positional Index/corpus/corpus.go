package corpus

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
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
				id:        id,
				file:      file,
				frequency: 1,
				positions: []int32{int32(position) + 1},
			})
			corpus.Put(w, Index{ Docs{docs}, 1})
		} else {
			documents := index.(Index)
			documents.totalFrequency++
			if !index.(Index).Contains(id) {
				documents.docs.Put(id, Doc{
					id: id,
					file:file,
					frequency: 1,
					positions:[]int32{int32(position) + 1},
				})
			} else {
				documents.updateDocument(id, position + 1)
			}
			corpus.Put(w, documents)
		}
	}
}

func (index *Index) updateDocument(id, position int) {
	document, _ := index.docs.Get(id)
	doc := document.(Doc)
	doc.frequency++
	doc.positions = append(doc.positions, int32(position))
	index.docs.Put(id, doc)
}

func splitRaw(s string) []string {
	return strings.Split(strings.Trim(s, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@"), " ")
}

func (corpus *Corpus) Print() {
	corpus.Each(func(key interface{}, value interface{}) {
		index := value.(Index)
		fmt.Printf("term: %s, total frequency: %d, posting list: \n",key.(string), index.totalFrequency)
		index.docs.Each(func(key interface{}, value interface{}) {
			fmt.Println(value.(Doc))
		})
	})
}
