package spimi

import (
	"fmt"
	"sync"
	"testing"
)

var spimi =  &SPIMI{
	inputDir:   "data",
	outputFile: "output/index.txt",
	blockSize:  5000,
	mutex: &sync.Mutex{},
	wg: &sync.WaitGroup{},

}

//func TestGenerateTokens(t *testing.T) {
//	fmt.Println(spimi.generateTokens())
//}
//
//func TestMakeBlocks(t *testing.T) {
//	tokenStream := spimi.generateTokens()
//	_ = spimi.makeBlocks(tokenStream)
//}

//func TestMergeBlocks(t *testing.T) {
//	tokenStream := spimi.generateTokens()
//	blocks := spimi.makeBlocks(tokenStream)
//	terms := getTerms(tokenStream)
//	spimi.mergeBlocks(terms, blocks)
//}
func TestSPIMI(t *testing.T) {
	fmt.Println(Spimi("data", "output/index.txt", 5000).FuzzySearch("world", 1))
}