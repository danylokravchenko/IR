package lemmatization

import (
	"../normalization"
	"bufio"
	"errors"
	"github.com/emirpasic/gods/maps/treemap"
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

var (
	ukranianConfig = Config{
		File:      "dictionary/lemmatization-ukr.txt",
		Separator: "	",
		Language:    "ukr",
	}
	russianConfig = Config{
		File:      "dictionary/lemmatization-rus.dat",
		Separator: "	",
		Language:    "rus",
	}
	englishConfig = Config{
		File:      "dictionary/lemmatization-ukr.txt",
		Separator: "	",
		Language:    "eng",
	}
)

// configure lemmatizator and read words from dictionary
func New(config Config) (*Lemmatizator, error) {

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
		line := strings.Split(normalizator.BasicNormalize(scanner.Text()), config.Separator)
		if len(line) >= 2 {
			if config.Language == "rus" {
				lemmatizator.tree.Put(line[0], line[1])
			} else {
				lemmatizator.tree.Put(line[1], line[0])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if lemmatizator.tree.Size() == 0 {
		return nil, errors.New("File is empty or format is wrong")
	}

	return lemmatizator, nil

}

func NewUkranianLemmatizator() (*Lemmatizator, error) {
	return New(ukranianConfig)
}

func NewRussianLemmatizator() (*Lemmatizator, error) {
	return New(ukranianConfig)
}
func NewEnglishLemmatizator() (*Lemmatizator, error) {
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
