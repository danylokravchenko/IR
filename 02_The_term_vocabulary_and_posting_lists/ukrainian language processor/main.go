package main

import (
	"./lemmatization"
	"./normalization"
	"./detector"
	"fmt"
	"log"
)

func main() {

	language := detector.DetectLanguage("абонент	абонентові	абонент	абонентом	абонент	абоненту.")
	fmt.Println(language)

	normalizator := normalization.New()

	ukrLemmatizator, err := lemmatization.NewUkranianLemmatizator()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ukrLemmatizator.GetLema(normalizator.Normalize("авАрі??!}Ях")))
	fmt.Println(ukrLemmatizator.Contains("аваріях"))

	lemmatizator, err := lemmatization.New(lemmatization.Config{
		File:      "dictionary/lemmatization-ukr.txt",
		Separator: "	",
		Language:    language,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(lemmatizator.Contains("аваріях"))

}
