package map_reduce

import (
	"fmt"
	"runtime"
	"testing"
)

func TestMapReduce(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Processing. Please wait....")

	// start the enumeration of files to be processed into a channel
	input := enumerateFiles("data")

	//TODO:// get length of input dir and pass as argument

	// this will start the map reduce work
	_ = mapReduce(mapper, reducer, inverter, input, 2)

	//fmt.Println(results)

	fmt.Println("Done!")
}
