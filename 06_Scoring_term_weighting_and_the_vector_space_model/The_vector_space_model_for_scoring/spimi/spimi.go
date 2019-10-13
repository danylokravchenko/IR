package spimi

import (
	. "../corpus"
	"bufio"
	"errors"
	"fmt"
	"github.com/emirpasic/gods/maps/hashmap"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
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
	inputDir      string
	outputFile    string
	tempBlockSize int
	termsInBlock  int
	docsNum       int
	corpus        *Corpus
	blockTree     *BlockTree
	mutex         *sync.Mutex
	wg            *sync.WaitGroup
}


func Spimi(inputDir, outputFile string, tempBlockSize, termsInBlock int) *BlockTree{

	spimi := &SPIMI{
		inputDir:      inputDir,
		outputFile:    outputFile,
		tempBlockSize: tempBlockSize,
		termsInBlock:  termsInBlock,
		mutex:  	   &sync.Mutex{},
		wg: 		   &sync.WaitGroup{},
	}
	tokenStream := spimi.generateTokens()
	blocks := spimi.makeTempBlocks(tokenStream)
	terms := getTerms(tokenStream)
	spimi.mergeTempBlocks(terms, blocks)

	spimi.createBlockStorage()

	return spimi.blockTree

}

func getTerms(tokens []Token) []string {

	res := make([]string, 0)

	for _, token := range tokens {
		res = append(res, token.Term)
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

	spimi.docsNum = len(files)

	spimi.wg.Add(len(files))

	for i, f := range files {
		// TODO: use concurrency to parse document
		// 1 gorutine per file
		go func(i int, filename string) {
			tokens, err := spimi.parseDocument(i, filename)
			if err != nil {
				log.Println(err)
				return
			}
			tokenStream = append(tokenStream, tokens...)
		}(i, spimi.inputDir +"/" + f.Name())

	}

	spimi.wg.Wait()

	return tokenStream

}


// Split string
func tokenize(text string) []string {
	//re := regexp.MustCompile(`(?:[A-Z]\.)+|\w+(?:[-']\w+)*|[-.(]+|\S\w*`)
	//return re.FindAllString(text, -1)
	tokens := strings.Split(strings.Trim(text, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@"), " ")
	for i, token := range tokens {
		tokens[i] = strings.Replace(
			strings.Replace(
			strings.Replace(
			strings.Replace(
			strings.Replace(
			strings.Replace(
			strings.Replace(
				token, ".", "", -1),
				",", "", -1),
				"'", "", -1),
				"?", "", -1),
				"!", "", -1),
				"\"", "", -1),
		";", "", -1)

	}
	return tokens
}

// Parse the doc in a stream of term-docId pairs which we call tokens
func (spimi *SPIMI) parseDocument(docID int, fileName string) ([]Token, error) {

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
				Term:     term,
				Position: pos + 1,
				DocID:    docID,
				File: fileName,
			})
		}
	}

	spimi.wg.Done()

	return tokens, nil

}

type blocks []string

// Create block channel of created files with tokens
func (spimi *SPIMI) makeTempBlocks(tokenStream []Token) blocks {

	blockID, begin, end := 0, 0, 0
	tokensCount := len(tokenStream)

	blocks := make(blocks, 0)

	for end + spimi.tempBlockSize < tokensCount {
		end += spimi.tempBlockSize
		blocks  = append(blocks, spimi.Invert(blockID, tokenStream[begin:end]))
		blockID++
		begin = end
	}

	blocks  = append(blocks, spimi.Invert(blockID, tokenStream[end:]))

	return blocks

}

// Create inverted Index
func (spimi *SPIMI) Invert(blockID int, tokens []Token) string{

	outputFile := fmt.Sprintf("temp/block%d.dat", blockID)

	c := NewCorpus()
	c.BuildIndexFromTokens(tokens)

	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	_, err = w.Write([]byte(c.ToGOB64()))
	if err != nil {
		log.Println(err)
	}
	w.Flush()

	return outputFile
}

func (spimi *SPIMI) mergeTempBlocks(terms []string, blocks blocks) {

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

		spimi.corpus.DocsNum = spimi.docsNum

		f.Close()
		//delete file
		err = os.Remove(b)
		if err != nil {
			log.Println(err)
		}
	}

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
			spimi.corpus.Put(term, Index{index.Docs, index.TotalFrequency, index.DocsNum, index.InverseDocumentFrequency})
			spimi.corpus.TermsNum++
		} else {
			documents := v.(Index)
			documents.TotalFrequency += index.TotalFrequency
			index.Docs.Each(func(key, value interface{}) {
				documents.DocsNum++
				docID := key.(int)
				doc := value.(Doc)
				if !documents.Contains(docID) {
					documents.Docs.Put(docID, doc)
				} else {
					documents.UpdateDocument(docID, doc.Positions)
				}
			})
			spimi.corpus.Put(term, documents)
		}
	})

}

func (spimi *SPIMI) createBlockStorage() {

	file, err := os.OpenFile(spimi.outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	spimi.corpus.CountInverseDocumentsFrequency()

	spimi.corpus.CountNormalizedDocumentsFrequency()

	keys := spimi.corpus.Keys()
	termsCount := len(keys)
	begin, end, blockID := 0,0,0

	spimi.blockTree = &BlockTree{hashmap.New(), spimi.corpus.Documents}

	spimi.wg.Add((termsCount/spimi.termsInBlock)-1)

	for end + spimi.termsInBlock < termsCount {
		end += spimi.termsInBlock
		terms := keys[begin:end]
		go func(blockID int, terms []interface{}) {
			block := spimi.createBlock(blockID, terms)
			for _, t := range terms {
				spimi.mutex.Lock()
				spimi.blockTree.Put(t, block)
				spimi.mutex.Unlock()
			}
			spimi.wg.Done()
		}(blockID, terms)

		blockID++
		begin = end
	}


	terms := keys[begin:]
	block := spimi.createBlock(blockID, terms)
	for _, t := range terms {
		spimi.mutex.Lock()
		spimi.blockTree.Put(t, block)
		spimi.mutex.Unlock()
	}

	spimi.wg.Wait()

	w := bufio.NewWriter(file)
	_, err = w.Write([]byte(spimi.blockTree.ToGOB64()))
	if err != nil {
		log.Println(err)
	}
	w.Flush()

}

func (spimi *SPIMI) createBlock(blockID int, terms []interface{}) string {

	outputFile := fmt.Sprintf("blocks/block%d.dat", blockID)

	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	tokens := make([]SerializedToken, 0)

	for _, t := range terms {
		idx, _ := spimi.corpus.Get(t)
		index := idx.(Index)
		docs := index.Docs.Values()
		documents := make([]SerializedDoc, 0)
		for _, d := range docs {
			doc := d.(Doc)
			documents = append(documents, SerializedDoc{
				Positions: doc.Positions,
				DocID:     doc.ID,
				File:      doc.File,
				Frequency: doc.Frequency,
				InverseDocumentFrequency: doc.InverseDocumentFrequency,
			})
		}
		tokens = append(tokens, SerializedToken{
			Term: t.(string),
			Docs: documents,
			TotalFrequency: index.TotalFrequency,
			InverseDocumentFrequency: index.InverseDocumentFrequency,
		})
	}

	sc := &SerializedCorpus{tokens}

	w := bufio.NewWriter(file)
	_, err = w.Write([]byte(sc.ToGOB64()))
	if err != nil {
		log.Println(err)
	}
	w.Flush()

	return outputFile

}