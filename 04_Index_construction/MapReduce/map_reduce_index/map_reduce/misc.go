package map_reduce

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"../corpus"
)

type WalkFunc func(path string, info os.FileInfo, err error) error

// get filenames in dir
func enumerateFiles(dirname string) chan interface{} {
	output := make(chan interface{})
	go func() {
		filepath.Walk(dirname, func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				output <- path
			}
			return nil
		})
		close(output)
	}()
	return output
}

// read file by line
func enumerateFile(filename string) chan string {
	output := make(chan string)
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			return
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')

			if err == io.EOF {
				break
			}

			// add each line to our enumeration channel
			output <- line
		}
		close(output)
	}()
	return output
}

// Split line
func tokenize(line string) []string {
	//re := regexp.MustCompile(`(?:[A-Z]\.)+|\w+(?:[-']\w+)*|[-.(]+|\S\w*`)
	//return re.FindAllString(text, -1)
	tokens := strings.Split(strings.Trim(line, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@"), " ")
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

func createSegmentFiles(tokens []corpus.Token, idx int) string {

	outputFile := fmt.Sprintf("output/segment%d.dat", idx)

	st := SerializedTokens{tokens}

	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	_, err = w.Write([]byte(st.ToGOB64()))
	if err != nil {
		log.Println(err)
	}
	w.Flush()

	return outputFile

}

type SerializedTokens struct {
	Tokens []corpus.Token
}

//// Serialize tokens
func (st SerializedTokens) ToGOB64() string {

	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(&st)
	if err != nil { fmt.Println(`failed gob Encode`, err) }

	return base64.StdEncoding.EncodeToString(b.Bytes())

}

// Deserialize tokens
// Go binary decoder
func FromGOB64(str string) []corpus.Token {

	st := &SerializedTokens{}

	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(`failed base64 Decode`, err);
	}

	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(st)
	if err != nil {
		fmt.Println(`failed gob Decode`, err);
	}

	return st.Tokens

}