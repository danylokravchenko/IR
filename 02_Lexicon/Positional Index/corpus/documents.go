package corpus

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
)

type Doc struct {
	id        int // int because map comparator is int type
	file      string
	frequency int32
	positions []int32
}

type Docs struct {
	*treemap.Map
}

// override toString
func (doc Doc) String() string {
	return fmt.Sprintf("{ID: %d, file: %s, frequency: %d, positions: %d", doc.id, doc.file, doc.frequency, doc.positions)
}

type Index struct {
	docs Docs //[]Doc
	totalFrequency int32
}

func (this Index) Contains(id int) bool {
	_, contains := this.docs.Get(id)
	return contains
}
