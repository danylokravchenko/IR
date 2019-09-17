package lemmatization

import (
	"bufio"
	"errors"
	"github.com/emirpasic/gods/maps/treemap"
	"os"
	"strings"
)

type Lemmatizer interface {
	New(dictionary string) *Lemmatizer
	GetLema(word string) string
	Contains(word string) bool
}

type UkranianLemmatizer struct {
	tree *treemap.Map
}

// read words from dictionary
func New(dictionary string) (*UkranianLemmatizer, error) {

	lemmatizer := &UkranianLemmatizer{
		treemap.NewWithStringComparator(),
	}
	file, err := os.Open(dictionary)
	if err != nil {
		return nil, errors.New("Error while opening a file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "	")
		if len(line) == 2 {
			lemmatizer.tree.Put(line[1], line[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lemmatizer, nil

}

// return lema if it exists or the same word
func (l *UkranianLemmatizer) GetLema(word string) string {
	if lema, ok := l.tree.Get(word); ok {
		return lema.(string)
	} else {
		return word
	}
}

// check if tree contains this key
func (l *UkranianLemmatizer) Contains(word string) bool {
	_, ok := l.tree.Get(word)
	return ok
}
