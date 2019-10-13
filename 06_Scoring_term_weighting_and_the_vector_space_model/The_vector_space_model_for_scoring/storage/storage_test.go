package storage

import (
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {
	//bt := InitStorage("/home/danil/Проекты/Go/Information_Retrieval/06_Scoring_term_weighting_and_the_vector_space_model/The_vector_space_model_for_scoring/spimi/data")
	//fmt.Println(ITFScore(bt, "What", "did"))
}

func TestCosineSimilarity(t *testing.T) {
	//doc1 := parseToTokens("what did they say")
	//doc2 := parseToTokens(`brother did often mean what they say ? Do you mean what you say now do do do?`)
	//fmt.Println(DotProduct2(doc1, doc2))
	//fmt.Println(EuclideanLength2(doc1))
	//fmt.Println(EuclideanLength2(doc2))
	//fmt.Println(EuclideanLength2(doc1)*EuclideanLength2(doc2))
	//fmt.Println(CosineSimilarity2(doc1, doc2))
	//fmt.Println("----")
}

func TestCosineScore(t *testing.T) {

	bt := InitStorage("/home/danil/Проекты/Go/Information_Retrieval/06_Scoring_term_weighting_and_the_vector_space_model/The_vector_space_model_for_scoring/spimi/data")
	fmt.Println(ITFScore(bt, "What", "did"))
	fmt.Println("----")
	fmt.Println(CosineScore(bt, `What did`, 10))

}