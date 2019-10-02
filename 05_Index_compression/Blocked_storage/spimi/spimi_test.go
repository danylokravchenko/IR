package spimi

import (
	"fmt"
	"sync"
	"testing"
)

var spimi =  &SPIMI{
	inputDir:      "data",
	outputFile:    "blocks/index.dat",
	tempBlockSize: 5000,
	termsInBlock:  4,
	mutex:         &sync.Mutex{},
	wg:            &sync.WaitGroup{},

}


func TestSPIMI(t *testing.T) {
	fmt.Println(Spimi("data", "blocks/index.dat", 5000, 4).Size())
}