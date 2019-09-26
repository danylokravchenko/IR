package spimi

import (
	"fmt"
	"testing"
)

var spimi =  &SPIMI{
	inputDir:   "data",
	outputFile: "index.txt",
	blockSize:  5000,
}

//func TestGenerateTokens(t *testing.T) {
//	_ = spimi.generateTokens()
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
	fmt.Println(NewSpimi("data", "index.txt", 5000).corpus.Keys())
}