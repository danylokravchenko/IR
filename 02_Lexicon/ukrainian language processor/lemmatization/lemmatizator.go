package lemmatization

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"../normalization"
	"os"
	"strings"
)

type ILemmatizator interface {
	New(dictionary string) *ILemmatizator
	GetLema(word string) string
	Contains(word string) bool
}

type Config struct {
	File string
	Separator string
	Language string
}

type Lemmatizator struct {
	tree *treemap.Map
	Config Config
}

var ukranianConfig = Config{
	File:      "dictionary/lemmatization-ukr.txt",
	Separator: "	",
	Language:    "ukr",
}

// configure lemmatizator and read words from dictionary
func New(config Config) (*Lemmatizator, error) {

	fmt.Println("Loading data...")

	lemmatizator := &Lemmatizator{
		treemap.NewWithStringComparator(),
		config,
	}

	normalizator := normalization.New()

	file, err := os.Open(config.File)
	if err != nil {
		return nil, errors.New("Error while opening a file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(normalizator.Normalize(scanner.Text()), config.Separator)
		if len(line) == 2 {
			lemmatizator.tree.Put(line[1], line[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if lemmatizator.tree.Size() == 0 {
		return nil, errors.New("File is empty or format is wrong")
	}

	fmt.Println("Data loaded")

	return lemmatizator, nil

}

func NewUkranianLemmatizator() (*Lemmatizator, error) {
	return New(ukranianConfig)
}

// return lema if it exists or the same word
func (l *Lemmatizator) GetLema(word string) string {
	if lema, ok := l.tree.Get(word); ok {
		return lema.(string)
	} else {
		return word
	}
}

// check if tree contains this key
func (l *Lemmatizator) Contains(word string) bool {
	_, ok := l.tree.Get(word)
	return ok
}
