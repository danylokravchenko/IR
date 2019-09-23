package automaton

type entry struct {
	idx, val int
}

// Sparse vector is a vector that avoid memory waisting and stores other values as zero
// Usually sparse vector are represented by a tuple (id, value) such as:
// ui=values[j] if id[j]=i; ui=0 otherwise (if i is not in id)
// And for example a dense vector (1, 2, 0, 0, 5, 0, 9, 0, 0)
// will be represented as {(0,1,4,6), (1, 2, 5, 9)}

// sparseVector is a crude implementation of a sparse vector for our needs
type sparseVector []*entry

// newSparseVector creates a new sparse vector with the initial values of the dense int slice given to it
func newSparseVector(values []int) sparseVector {
	v := make(sparseVector, len(values))

	for i := 0; i < len(values); i++ {
		v[i] = &entry{i, values[i]}
	}

	return v
}

// append appends another sparse vector entry with the given index and value. NOTE: We do not check
// that an entry with the same index is present in the vector
func (v sparseVector) append(index, value int) sparseVector {
	return append(v, &entry{idx: index, val: value})
}
