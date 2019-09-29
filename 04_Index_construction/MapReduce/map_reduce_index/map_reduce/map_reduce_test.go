package map_reduce

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"testing"
	"../corpus"
)

func getFilesLength(path string) int {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	return len(files)
}

func TestMapReduce(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Processing. Please wait....")

	//TODO: I don't know why, but big files rape this program
	// RAM is fully loaded

	// start the enumeration of files to be processed into a channel
	input := enumerateFiles("data")

	// get amount of files in dir, just a `hack`
	length := getFilesLength("data")

	// this will start the map reduce work
	c := mapReduce(mapper, reducer, inverter, input, length)

	fmt.Println(c.(*corpus.Corpus).FuzzySearch("world", 1))

	fmt.Println("Done!")
}
