package storage

import (
	"../spimi"
	"../corpus"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	outputFile = "blocks/index.dat"
	tempBlockSize = 5000
	termsInBlock = 4
)

func InitStorage(inputDir string) *corpus.BlockTree {

	var bt *corpus.BlockTree

	if !fileExists(outputFile) {
		bt = spimi.Spimi(inputDir, outputFile, tempBlockSize, termsInBlock)
	}

	bt = loadBTree(outputFile)

	return bt

}

func DeserializeBlock(path string) *corpus.SerializedCorpus{
	f , err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

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
				fmt.Println("Error reading file", err);
				os.Exit(1)
			}
		}
	}

	return corpus.SerializedCorpusFromBlock(string(data))

}

func DeserializeTerm(term, path string) corpus.SerializedToken {
	return DeserializeBlock(path).Filter(func(token corpus.SerializedToken) bool{
		return token.Term == term
	}).Tokens[0]
}

func fileExists(path string) bool {
	// detect if file exists
	var _, err = os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return true

}

func loadBTree(path string) *corpus.BlockTree {

	f , err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

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
				fmt.Println("Error reading file", err);
				os.Exit(1)
			}
		}
	}
	bt := corpus.BlockTreeFromGOB64(string(data))

	return bt

}
