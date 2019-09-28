package map_reduce

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"testing"
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

	// start the enumeration of files to be processed into a channel
	input := enumerateFiles("data")

	// get amount of files in dir, just a `hack`
	length := getFilesLength("data")

	// this will start the map reduce work
	_ = mapReduce(mapper, reducer, inverter, input, length)

	fmt.Println("Done!")
}
