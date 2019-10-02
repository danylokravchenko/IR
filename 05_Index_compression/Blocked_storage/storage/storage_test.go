package storage

import (
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {
	bt := InitStorage("/home/danil/Проекты/Go/Information Retrieval/05_Index_compression/Blocked_storage/spimi/data")
	term := "world"
	if block, ok := bt.Get(term); ok {
		fmt.Println(DeserializeBlock(term, block.(string)))
	}
}
