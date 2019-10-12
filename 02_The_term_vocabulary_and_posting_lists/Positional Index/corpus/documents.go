package corpus

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
)

type Doc struct {
	ID        int // int because map comparator is int type
	File      string
	Frequency int32
	Positions []int32
}

type Docs struct {
	*treemap.Map
}

// override toString
func (doc Doc) String() string {
	return fmt.Sprintf("{ID: %d, File: %s, NormalizedFrequency: %d, Positions: %d", doc.ID, doc.File, doc.Frequency, doc.Positions)
}

type Index struct {
	Docs           Docs //[]Doc
	TotalFrequency int32
}

func (this Index) Contains(id int) bool {
	_, contains := this.Docs.Get(id)
	return contains
}
