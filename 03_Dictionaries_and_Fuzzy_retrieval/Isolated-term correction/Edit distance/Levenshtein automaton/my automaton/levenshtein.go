package my_automaton


// Sparse automaton is implementation of Fast search Levenshtein automaton
// str - string to check
// max - max number of edits (edit distance)
type SparseAutomaton struct {
	str string
	max int
}


// Return new sparse initialized automaton
func NewSparseAutomaton(word string, maxEdits int) *SparseAutomaton {
	return &SparseAutomaton{
		str: word,
		max: maxEdits,
	}
}


// Initialize the automaton's state and return sparseVector for the iteration over next steps
func (a *SparseAutomaton) Start() sparseVector {
	values := make([]int, a.max + 1)

	for i := range values {
		values[i] = i
	}

	return newSparseVector(values)
}


// helper to find minimal value of 2 given
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}


// Step returns the next state of the automaton given a previous state and a character to check
func (a *SparseAutomaton) Step(state sparseVector, c byte) sparseVector {

	newVector := make(sparseVector, 0)

	if len(state) > 0 && state[0].idx == 0 && state[0].val < a.max {
		newVector = newVector.append(0, state[0].val+1)
	}

	for i, entry := range state {

		if entry.idx == len(a.str) {
			break
		}

		cost := 0
		if a.str[entry.idx] != c {
			cost = 1
		}

		val := state[i].val + cost

		if len(newVector) != 0 && newVector[len(newVector)-1].idx == entry.idx {
			val = min(val, newVector[len(newVector)-1].val+1)
		}

		if len(state) > i+1 && state[i+1].idx == entry.idx+1 {
			val = min(val, state[i+1].val+1)
		}

		if val <= a.max {
			newVector = newVector.append(entry.idx+1, val)
		}
	}

	return newVector
}


// IsMatch returns true if the current state vector represents a string that is within the max
// edit distance from the initial automaton string
func (a *SparseAutomaton) IsMatch(v sparseVector) bool {
	return len(v) != 0 && v[len(v)-1].idx == len(a.str)
}


// CanMatch returns true if there is a possibility that feeding the automaton with more steps will
// yield a match. Once CanMatch is false there is no point in continuing iteration
func (a *SparseAutomaton) CanMatch(v sparseVector) bool {
	return len(v) > 0
}


func (a *SparseAutomaton) Transitions(v sparseVector) []byte {

	set := map[byte]struct{}{}
	for _, entry := range v {
		if entry.idx < len(a.str) {
			set[a.str[entry.idx]] = struct{}{}
		}
	}

	ret := make([]byte, 0, len(set))
	for b, _ := range set {
		ret = append(ret, b)
	}

	return ret

}
