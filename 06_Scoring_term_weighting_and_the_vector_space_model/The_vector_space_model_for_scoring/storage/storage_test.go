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

	bt := InitStorage("/home/danil/Проекты/Go/Information_Retrieval/06_Scoring_term_weighting_and_the_vector_space_model/The_vector_space_model_for_scoring/spimi/data")

	fmt.Println(CosineScore(bt, `what did?`, 1))

}