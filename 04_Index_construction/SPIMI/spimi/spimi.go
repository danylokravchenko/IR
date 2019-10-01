package spimi

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//Page 73:
//SPIMI-Invert(token_stream)
//1  output_file = NewFile()
//2  dictionary = NewHash()
//3  while (free memory available)
//4  do token <- next(token_stream)
//5    if term(token) not in dictionary
//6      then postings_list = AddToDictionary(dictionary,term(token))
//7      else postings_list = GetPostingsList(dictionary,term(token))
//8    if full(postings_list)
//9      then postings_list = DoublePostingsList(dictionary,term(token))
//10    AddToPostingsList(postings_list, docID(token))
//11  sorted_terms <- SortTerms(dictionary)
//12  WriteBlockToDisk(sorted_terms,dictionary,output_file)
//13  return output_file
//Each call of SPIMI-Invert writes a block to disk.
//The index of the block is its dictionary and the postings_list.

type SPIMI struct {
	inputDir string
	outputFile string
	blockSize int
	corpus *Corpus
}

func NewSpimi(inputDir, outputFile string, blockSize int) *SPIMI{

	spimi := &SPIMI{
		inputDir:   inputDir,
		outputFile: outputFile,
		blockSize:  blockSize,
	}
	tokenStream := spimi.generateTokens()
	blocks := spimi.makeBlocks(tokenStream)
	terms := getTerms(tokenStream)
	spimi.mergeBlocks(terms, blocks)

	return spimi

}

func getTerms(tokens []Token) []string {
	res := make([]string, 0)
	for _, token := range tokens {
		res = append(res, token.term)
	}
	return res
}

// Generate tokens from files in data dir
// TODO: Maybe good idea is to return chanel, so program could run forward while this method will parse files in dir
func (spimi *SPIMI) generateTokens() []Token {

	tokenStream := make([]Token, 0)

	files, err := ioutil.ReadDir(spimi.inputDir)
	if err != nil {
		// TODO: Warning!!!!
		log.Fatal(err)
	}

	for i, f := range files {
		// TODO: use concurrency to parse document
		// 1 gorutine per file
		tokens, err := parseDocument(i, spimi.inputDir +"/" + f.Name())
		if err != nil {
			log.Println(err)
			continue
		}
		tokenStream = append(tokenStream, tokens...)
	}

	return tokenStream

}


// Split string
func tokenize(text string) []string {
	//re := regexp.MustCompile(`(?:[A-Z]\.)+|\w+(?:[-']\w+)*|[-.(]+|\S\w*`)
	//return re.FindAllString(text, -1)
	tokens := strings.Split(strings.Trim(text, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@"), " ")
	for i, token := range tokens {
		if strings.HasSuffix(token, ".") || strings.HasSuffix(token, ",") || strings.HasSuffix(token, ";") {
			tokens[i] = token[:len(token)-1]
		}
	}
	return tokens
}

type Token struct {
	term string
	position int
	docID int
	file string
}

// Parse the doc in a stream of term-docId pairs which we call tokens
func parseDocument(docID int, fileName string) ([]Token, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return []Token{}, errors.New(fmt.Sprintf("Cannot open file: %v", err))
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	tokens := make([]Token, 0)

	for scanner.Scan() {
		terms := tokenize(scanner.Text())
		for pos, term := range terms {
			tokens = append(tokens, Token{
				term:     term,
				position: pos,
				docID:    docID,
				file: fileName,
			})
		}
	}

	return tokens, nil

}

type blocks []string

// Create block channel of created files with tokens
func (spimi *SPIMI) makeBlocks(tokenStream []Token) blocks {

	blockID, begin, end := 0, 0, 0
	tokensCount := len(tokenStream)

	blocks := make(blocks, 0)

	for end + spimi.blockSize < tokensCount {
		end += spimi.blockSize
		blocks  = append(blocks, spimi.Invert(blockID, tokenStream[begin:end]))
		blockID++
		begin = end
	}

	blocks  = append(blocks, spimi.Invert(blockID, tokenStream[end:]))

	return blocks

}

func (spimi *SPIMI) mergeBlocks(terms []string, blocks blocks) {

	file, err := os.OpenFile(spimi.outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	for _, b := range blocks {
		f , err := os.Open(b)
		if err != nil {
			log.Println(err)
		}

		stat, err := f.Stat()
		if err != nil {
			log.Println(err)
		}
		data := make([]byte, stat.Size())
		for {
			_, err = f.Read(data)
			if err != nil {
				if err == io.EOF {
					break // end of the file
				} else {
					fmt.Println("Error reading file");
					os.Exit(1)
				}
			}
		}
		corpus := FromGOB64(string(data))
		spimi.IntersectCorpuses(corpus)

		f.Close()
		//delete file
		err = os.Remove(b)
		if err != nil {
			log.Println(err)
		}
	}


	w := bufio.NewWriter(file)
	_, err = w.Write([]byte(spimi.corpus.ToGOB64()))
	if err != nil {
		log.Println(err)
	}
	w.Flush()

}


func (spimi *SPIMI) IntersectCorpuses(c *Corpus) {

	if spimi.corpus == nil {
		spimi.corpus = c
		return
	}

	c.Each(func(key, value interface{}) {
		term := key.(string)
		index := value.(Index)
		if v, ok := spimi.corpus.Get(term); !ok {
			spimi.corpus.Put(term, Index{index.Docs, index.TotalFrequency})
		} else {
			documents := v.(Index)
			documents.TotalFrequency += index.TotalFrequency
			index.Docs.Each(func(key, value interface{}) {
				docID := key.(int)
				doc := value.(Doc)
				if !documents.Contains(docID) {
					documents.Docs.Put(docID, doc)
				} else {
					documents.updateDocument(docID, doc.Positions)
				}
			})
		}
	})

}


// Create inverted Index
func (spimi *SPIMI) Invert(blockID int, tokens []Token) string{

	outputFile := fmt.Sprintf("output/block%d.dat", blockID)

	corpus := &Corpus {treemap.NewWithStringComparator()}

	for _, token := range tokens {
		if index, ok := corpus.Get(token.term); !ok {
			docs := treemap.NewWithIntComparator()
			docs.Put(token.docID, Doc{
				ID:        token.docID,
				File:      token.file,
				Frequency: 1,
				Positions: []int{token.position + 1},
			})
			corpus.Put(token.term, Index{Docs{docs}, 1})
		} else {
			documents := index.(Index)
			documents.TotalFrequency++
			if !index.(Index).Contains(token.docID) {
				documents.Docs.Put(token.docID, Doc{
					ID:        token.docID,
					File:      token.file,
					Frequency: 1,
					Positions: []int{token.position + 1},
				})
			} else {
				documents.updateDocument(token.docID, []int{token.position+1})
			}
		}
	}

	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	_, err = w.Write([]byte(corpus.ToGOB64()))
	if err != nil {
		log.Println(err)
	}
	w.Flush()

	return outputFile
}
