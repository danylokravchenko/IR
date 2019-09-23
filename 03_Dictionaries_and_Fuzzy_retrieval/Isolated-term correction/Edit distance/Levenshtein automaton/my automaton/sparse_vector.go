package my_automaton


// Sparse vector is a vector that avoid memory waisting and stores other values as zero
// Usually sparse vector are represented by a tuple (id, value) such as:
// ui=values[j] if id[j]=i; ui=0 otherwise (if i is not in id)
// And for example a dense vector (1, 2, 0, 0, 5, 0, 9, 0, 0)
// will be represented as {(0,1,4,6), (1, 2, 5, 9)}
type entry struct {
	idx int
	val int
}

type sparseVector []*entry


// Creates a new sparse vector from the given dense int slice
func newSparseVector(values []int) sparseVector {

	vector := make (sparseVector, len(values))

	for i, v := range values {
		vector[i] = &entry{ idx: i,	val: v }
	}

	return vector

}


// Append another entry
func (v sparseVector) append(idx, val int) sparseVector {
	return append(v, &entry{ idx: idx, val: val })
}