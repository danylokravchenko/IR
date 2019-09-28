package map_reduce

import (
	"../corpus"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Flow:
// a) Parser:
// 1) parse documents and get tokens
// 2) map function => sent parsed tokens into reducer
// 3) divide parsed tokens into diapasons of different first letters
// 4) reduce function => save results into segment files
// b) Inverter:
// 5) inverter => create 1 big inverted index from segment files

// MapperCollector is a channel that collects the output from mapper tasks
type MapperCollector chan chan interface{}

// MapperFunc is a function that performs the mapping part of the MapReduce job
type MapperFunc func(interface{}, chan interface{})

// ReducerFunc is a function that performs the reduce part of the MapReduce job
type ReducerFunc func(chan interface{}, chan interface{})

// InverterFunc is a function that performs the invert part of the MapReduce job
type InverterFunc func(chan interface{}, chan interface{}, int)

func mapper(filename interface{}, output chan interface{}) {

	tokens := make([]corpus.Token, 0)

	docID := 0

	// start the enumeration of each line in the file
	for line := range enumerateFile(filename.(string)) {

		terms := tokenize(line)

		for pos, term := range terms {
			tokens = append(tokens, corpus.Token{
				Term:     term,
				Position: pos + 1,
				DocID:    docID,
				File: 	  filename.(string),
			})
		}

		docID++
	}

	output <- tokens

}

func reducer(input, output chan interface{}) {
	//results := make(blocks, 0)
	results := make(chan string)
	idx := 0
	for tokens := range input {
		go func(tokens []corpus.Token, idx int) {
			results <- createSegmentFiles(tokens, idx)
		}(tokens.([]corpus.Token), idx)
		// write tokens into segment file and save filename
		//results = append(results, createSegmentFiles(tokens.([]corpus.Token), idx))
		idx++
	}
	for res := range results {
		output <- res
	}

	close(results)

}

// open segment files and create Corpus
func inverter (input, output chan interface{}, length int) {
	// corpus
	result := make(chan []corpus.Token)

	wg := &sync.WaitGroup{}

	counter := 0

	for  {

		if counter == length {
			break
		}

		b := <- input

		wg.Add(1)

		go func(filename string) {

			f, err := os.Open(filename)
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
						fmt.Println("Error reading file");
						os.Exit(1)
					}
				}
			}

			result <- FromGOB64(string(data))
			wg.Done()

		}(b.(string))

		counter++
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	tokens := make([]corpus.Token, 0)

	for res := range result {
		tokens = append(tokens, res...)
	}

	c := corpus.NewCorpus(2)
	c.BuildIndexFromTokens(tokens)
	output <- c
}

// The mapperDispatcher function is responsible to listen on the data channel that receives each filename
// to be processed and invoke a mapper for each file, pushing the output of the job into a MapperCollector
func mapperDispatcher(mapper MapperFunc, input chan interface{}, collector MapperCollector) {
	for item := range input {
		taskOutput := make(chan interface{})
		go mapper(item, taskOutput)
		collector <- taskOutput
	}
	close(collector)
}

// The reducerDispatcher function is responsible to listen on the collector channel
// and push each item as the data for the Reducer task.
func reducerDispatcher(collector MapperCollector, reducerInput chan interface{}) {
	for output := range collector {
		reducerInput <- <-output
	}
	close(reducerInput)
}

const (
	MaxWorkers = 10
)

func mapReduce(mapper MapperFunc, reducer ReducerFunc, inverter InverterFunc, input chan interface{}, length int) interface{} {

	reducerInput := make(chan interface{}, length)
	reducerOutput := make(chan interface{}, length)
	inverterOutput := make(chan interface{})
	mapperCollector := make(MapperCollector, MaxWorkers)

	go reducer(reducerInput, reducerOutput)
	go reducerDispatcher(mapperCollector, reducerInput)
	go mapperDispatcher(mapper, input, mapperCollector)
	go inverter(reducerOutput, inverterOutput, length)

	return <-inverterOutput

}