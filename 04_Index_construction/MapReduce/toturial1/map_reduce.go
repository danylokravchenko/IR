package toturial1

import (
	"fmt"
	"sync"
)

// mapper receives a channel of strings and counts the occurrence of each unique word read from this channel.
// It sends the resulting map to the output channel.
func mapper(in <-chan string, out chan<- map[string]int) {
	count := map[string]int{}
	for word := range in {
		count[word] = count[word] + 1
	}
	out <- count
	close(out)
}

// reducer receives a channel of ints and adds up all ints until the channel is closed.
// Then it divides through the number of received ints to calculate the average.
func reducer(in <-chan int, out chan<- float32) {
	sum, count := 0, 0
	for n := range in {
		sum += n
		count++
	}
	out <- float32(sum) / float32(count)
	close(out)
}

// inputDistributor receives three output channels and sends each of them some data.
func inputReader(out [3]chan<- string) {
	// "Read" some data.
	input := [][]string{
		{"noun", "verb", "verb", "noun", "noun"},
		{"verb", "verb", "verb", "noun", "noun", "verb"},
		{"noun", "noun", "verb", "noun"},
	}

	for i := range out {
		go func(ch chan<- string, word []string) {
			for _, w := range word {
				ch <- w
			}
			close(ch)
		}(out[i], input[i])
	}
}

// shuffler gets a list of data channels containing key/value pairs like
// "noun: 5, verb: 4". For each "noun" key, it sends the corresponding value
// to out[0], and for each "verb" key to out[1].
// The data channles are multiplexed into one, based on the `merge` function
// from the [Pipelines article](https://blog.golang.org/pipelines) of the
// Go Blog.
func shuffler(in []<-chan map[string]int, out [2]chan<- int) {
	var wg sync.WaitGroup
	wg.Add(len(in))
	for _, ch := range in {
		go func(c <-chan map[string]int) {
			for m := range c {
				nc, ok := m["noun"]
				if ok {
					out[0] <- nc
				}
				vc, ok := m["verb"]
				if ok {
					out[1] <- vc
				}
			}
			wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out[0])
		close(out[1])
	}()
}

// outputWriter starts a goroutine for each data channel and writes out
// the averages that it receives from each channel.
func outputWriter(in []<-chan float32) {
	var wg sync.WaitGroup
	wg.Add(len(in))
	// `out[0]` contains the nouns, `out[1]` the verbs.
	name := []string{"noun", "verb"}
	for i := 0; i < len(in); i++ {
		go func(n int, c <-chan float32) {
			for avg := range c {
				fmt.Printf("Average number of %ss per data text: %f\n", name[n], avg)
			}
			wg.Done()
		}(i, in[i])
	}
	wg.Wait()
}

func main() {
	// Set up all channels used for passing data between the workers.
	//
	// I could have used loops instead, to create arrays or
	// slices of channels. Apparently, copy/paste has won.
	size := 10
	text1 := make(chan string, size)
	text2 := make(chan string, size)
	text3 := make(chan string, size)
	map1 := make(chan map[string]int, size)
	map2 := make(chan map[string]int, size)
	map3 := make(chan map[string]int, size)
	reduce1 := make(chan int, size)
	reduce2 := make(chan int, size)
	avg1 := make(chan float32, size)
	avg2 := make(chan float32, size)

	// Start all workers in separate goroutines, chained together via channels.
	go inputReader([3]chan<- string{text1, text2, text3})
	go mapper(text1, map1)
	go mapper(text2, map2)
	go mapper(text3, map3)
	go shuffler([]<-chan map[string]int{map1, map2, map3}, [2]chan<- int{reduce1, reduce2})
	go reducer(reduce1, avg1)
	go reducer(reduce2, avg2)

	// The outputWriter runs in the main thread.
	outputWriter([]<-chan float32{avg1, avg2})
}
